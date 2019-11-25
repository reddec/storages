package filestorage

import (
	"github.com/pkg/errors"
	"github.com/reddec/storages"
	"io/ioutil"
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

func (ds *flatStorage) Namespace(name []byte) (storages.Storage, error) {
	dirName := string(name)
	if strings.ContainsRune(dirName, os.PathSeparator) {
		return nil, errors.New(errWithPathSeparator)
	}
	subLocation := ds.namespacePath(dirName)
	err := os.MkdirAll(subLocation, 0755)
	if err != nil {
		return nil, err
	}
	return NewFlat(dirName), nil
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
	err := os.MkdirAll(ds.fileNamePath(""), filePermission)
	if err != nil {
		return errors.Wrap(err, "create dir")
	}
	return filepath.Walk(ds.fileNamePath(""), func(path string, info os.FileInfo, err error) error {
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
	err := os.MkdirAll(ds.namespacePath(""), filePermission)
	if err != nil {
		return errors.Wrap(err, "create dir")
	}
	return filepath.Walk(ds.namespacePath(""), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		return handler([]byte(filepath.Base(path)))
	})
}

func (ds *flatStorage) Close() error { return nil } // NOP

func (ds *flatStorage) fileNamePath(name string) string {
	return filepath.Join(ds.location, "data", name)
}
func (ds *flatStorage) namespacePath(name string) string {
	return filepath.Join(ds.location, "namespace", name)
}
