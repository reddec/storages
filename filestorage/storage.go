package filestorage

import (
	"crypto"
	_ "crypto/sha256" // load for default
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/reddec/chop-text"
	"github.com/reddec/storages"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

const (
	metaDataFileSuffix = ".meta.json"
	chopSlices         = 4
	chopSize           = 4
	filePermission     = 0755
)

type dirStorage struct {
	location string
	chopper  chop.Chopper
	lock     sync.RWMutex
	hash     crypto.Hash
}

type metaInfo struct {
	Key []byte `json:"key"`
}

func (ds *dirStorage) Put(key []byte, data []byte) error {
	ds.lock.Lock()
	defer ds.lock.Unlock()
	targetFile := ds.getTargetFile(key)
	baseDir := path.Dir(targetFile)
	err := os.MkdirAll(baseDir, filePermission)
	if err != nil {
		return errors.Wrap(err, "create dir")
	}
	err = ioutil.WriteFile(targetFile, data, filePermission)
	if err != nil {
		return errors.Wrap(err, "put data to "+targetFile)
	}
	metaData, err := json.MarshalIndent(metaInfo{Key: key}, "", "  ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(ds.getMetaFileOfTarget(targetFile), metaData, filePermission)
	return errors.Wrap(err, "write meta data")
}

func (ds *dirStorage) Get(key []byte) ([]byte, error) {
	ds.lock.RLock()
	defer ds.lock.RUnlock()
	targetFile := ds.getTargetFile(key)
	data, err := ioutil.ReadFile(targetFile)
	if os.IsNotExist(err) {
		return nil, os.ErrNotExist
	}
	return data, errors.Wrap(err, "read key")
}

func (ds *dirStorage) Del(key []byte) error {
	ds.lock.Lock()
	defer ds.lock.Unlock()
	targetFile := ds.getTargetFile(key)
	metaFile := ds.getMetaFileOfTarget(targetFile)

	err := os.RemoveAll(metaFile)
	if err != nil {
		return errors.Wrap(err, "remove meta file")
	}

	err = os.RemoveAll(targetFile)
	return errors.Wrap(err, "remove data file")
}

func (ds *dirStorage) Keys(handler func(key []byte) error) error {
	ds.lock.RLock()
	defer ds.lock.RUnlock()
	return filepath.Walk(ds.location, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(info.Name(), metaDataFileSuffix) {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return errors.Wrap(err, "read meta file")
			}
			var meta metaInfo
			err = json.Unmarshal(data, &meta)
			if err != nil {
				return errors.Wrap(err, "parse meta data")
			}
			return handler(meta.Key)
		}
		return nil
	})
}

func (ds *dirStorage) Close() error { return nil } // NOP

func (ds *dirStorage) getTargetFile(key []byte) (string) {
	hash := hex.EncodeToString(ds.hash.New().Sum(key))
	return path.Join(ds.location, ds.chopper.Chop(hash))
}

func (ds *dirStorage) getMetaFileOfTarget(path string) string {
	return path + metaDataFileSuffix
}

// Same as NewHash but with SHA256 by default
func NewDefault(location string) storages.Storage {
	return New(location, crypto.SHA256)
}

// New file storage where each item stores in single file with path based on hashed key. Additionally near each file, meta file also generates
func New(location string, hash crypto.Hash) storages.Storage {
	return &dirStorage{
		chopper:  chop.Chopper{Sep: os.PathSeparator, Slices: chopSlices, SliceSize: chopSize},
		location: location,
		hash:     hash,
	}
}
