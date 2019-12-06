package filestorage

import (
	"github.com/pkg/errors"
	"github.com/reddec/storages"
	"github.com/reddec/storages/std"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	errWithPathSeparator = "name contains path separator"
)

func NewFlat(location string) storages.NamespacedStorage {
	return &flatStorage{
		location: location,
	}
}

type flatStorage struct {
	location string
	lock     sync.RWMutex
}

func (ds *flatStorage) DelNamespace(name []byte) error {
	dirName := string(name)
	if strings.ContainsRune(dirName, os.PathSeparator) {
		return errors.New(errWithPathSeparator)
	}
	return os.RemoveAll(filepath.Join(ds.location, dirName))
}

func (ds *flatStorage) Namespace(name []byte) (storages.Storage, error) {
	dirName := string(name)
	if strings.ContainsRune(dirName, os.PathSeparator) {
		return nil, errors.New(errWithPathSeparator)
	}
	subLocation := filepath.Join(ds.location, dirName)
	err := os.MkdirAll(subLocation, 0755)
	if err != nil {
		return nil, err
	}
	return NewFlat(subLocation), nil
}

func (ds *flatStorage) Put(key []byte, data []byte) error {
	ds.lock.Lock()
	defer ds.lock.Unlock()
	fileName := string(key)
	if strings.ContainsRune(fileName, os.PathSeparator) {
		return errors.New(errWithPathSeparator)
	}
	targetFile := ds.fileNamePath(fileName)
	err := os.MkdirAll(filepath.Dir(targetFile), filePermission)
	if err != nil {
		return errors.Wrap(err, "create dir")
	}
	err = ioutil.WriteFile(targetFile, data, filePermission)
	if err != nil {
		return errors.Wrap(err, "put data to "+targetFile)
	}
	return nil
}

func (ds *flatStorage) Get(key []byte) ([]byte, error) {
	ds.lock.RLock()
	defer ds.lock.RUnlock()
	fileName := string(key)
	if strings.ContainsRune(fileName, os.PathSeparator) {
		return nil, errors.New(errWithPathSeparator)
	}
	targetFile := ds.fileNamePath(fileName)
	data, err := ioutil.ReadFile(targetFile)
	if os.IsNotExist(err) {
		return nil, os.ErrNotExist
	}
	return data, errors.Wrap(err, "read key")
}

func (ds *flatStorage) Del(key []byte) error {
	ds.lock.Lock()
	defer ds.lock.Unlock()
	fileName := string(key)
	if strings.ContainsRune(fileName, os.PathSeparator) {
		return errors.New(errWithPathSeparator)
	}
	targetFile := ds.fileNamePath(fileName)
	err := os.RemoveAll(targetFile)
	return errors.Wrap(err, "remove file")
}

func (ds *flatStorage) Keys(handler func(key []byte) error) error {
	ds.lock.RLock()
	defer ds.lock.RUnlock()
	err := os.MkdirAll(ds.location, filePermission)
	if err != nil {
		return errors.Wrap(err, "create dir")
	}
	return filepath.Walk(ds.location, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		return handler([]byte(filepath.Base(path)))
	})
}

func (ds *flatStorage) Namespaces(handler func(name []byte) error) error {
	ds.lock.RLock()
	defer ds.lock.RUnlock()
	return filepath.Walk(ds.location, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() || path == ds.location {
			return nil
		}
		return handler([]byte(filepath.Base(path)))
	})
}

func (ds *flatStorage) Close() error { return nil } // NOP

func (ds *flatStorage) fileNamePath(name string) string {
	return filepath.Join(ds.location, name)
}

func init() {
	std.RegisterWithMapper("file", func(url *url.URL) (storage storages.Storage, e error) {
		return NewDefault(filepath.Join(url.Host, url.Path)), nil
	})
	std.RegisterWithMapper("file+flat", func(url *url.URL) (storage storages.Storage, e error) {
		return NewFlat(filepath.Join(url.Host, url.Path)), nil
	})
}
