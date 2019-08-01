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

func NewFlat(location string) storages.Storage {
	return &flatStorage{
		location: location,
	}
}

type flatStorage struct {
	location string
	lock     sync.RWMutex
}

func (ds *flatStorage) Put(key []byte, data []byte) error {
	ds.lock.Lock()
	defer ds.lock.Unlock()
	fileName := string(key)
	if strings.ContainsRune(fileName, os.PathSeparator) {
		return errors.New(errWithPathSeparator)
	}
	err := os.MkdirAll(ds.location, filePermission)
	if err != nil {
		return errors.Wrap(err, "create dir")
	}
	targetFile := filepath.Join(ds.location, fileName)
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
	targetFile := filepath.Join(ds.location, fileName)
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
	targetFile := filepath.Join(ds.location, fileName)
	err := os.RemoveAll(targetFile)
	return errors.Wrap(err, "remove file")
}

func (ds *flatStorage) Keys(handler func(key []byte) error) error {
	ds.lock.RLock()
	defer ds.lock.RUnlock()
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

func (ds *flatStorage) Close() error { return nil } // NOP
