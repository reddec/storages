package indexed

import (
	"bytes"
	"encoding/gob"
	"github.com/reddec/storages"
	"os"
	"sync"
)

// Creates new unique index backed by dedicated storage. Multiple links
// with a same secondary key will overwrite each other.
func NewUniqueIndex(storage storages.Storage) Index {
	return &uniqueIndex{index: storage}
}

// Creates new non-unique index backed by dedicated storage. Multiple links
// with a same secondary key will be packed to array (encoded by GOB, but it could be changed).
// Important! Due to process of appending/modify arrays is in memory, user should be aware of
// badly distributed secondary keys (where there are a lot of same secondary keys).
func NewIndex(storage storages.Storage) Index {
	return &multiIndex{index: storage}
}

// Index for secondary keys
type Index interface {
	// Link secondary key to primary key (like Email -> User ID)
	Link(primaryKey, secondaryKey []byte) error
	// Unlink secondary key
	Unlink(primaryKey, secondaryKey []byte) error
	// Find primary keys by secondary key
	Find(secondaryKey []byte) ([][]byte, error)
	// Iterate over entries in index. Order depends of underlying storage
	Iterate(handler func(primaryKey, secondaryKey []byte) error) error
}

type uniqueIndex struct {
	index storages.Storage
}

func (uq *uniqueIndex) Find(secondaryKey []byte) ([][]byte, error) {
	pk, err := uq.index.Get(secondaryKey)
	if err != nil {
		return nil, err
	}
	return [][]byte{pk}, nil
}

func (uq *uniqueIndex) Unlink(primaryKey, secondaryKey []byte) error {
	return uq.index.Del(secondaryKey)
}

func (uq *uniqueIndex) Link(primaryKey, secondaryKey []byte) error {
	return uq.index.Put(secondaryKey, primaryKey)
}

func (uq *uniqueIndex) Iterate(handler func(primaryKey, secondaryKey []byte) error) error {
	return uq.index.Keys(func(secondaryKey []byte) error {
		primaryKey, err := uq.index.Get(secondaryKey)
		if err != nil {
			return err
		}
		return handler(primaryKey, secondaryKey)
	})
}

type multiIndex struct {
	index storages.Storage
	lock  sync.RWMutex
}

func (mi *multiIndex) Find(secondaryKey []byte) ([][]byte, error) {
	mi.lock.RLock()
	defer mi.lock.RUnlock()
	return mi.getPrimaryKeys(secondaryKey)
}

func (mi *multiIndex) Unlink(primaryKey, secondaryKey []byte) error {
	mi.lock.Lock()
	defer mi.lock.Unlock()
	primaryKeys, err := mi.getPrimaryKeys(secondaryKey)
	if err != nil {
		return err
	}
	if len(primaryKeys) == 0 {
		return nil
	}
	var cp = make([][]byte, len(primaryKeys))
	var j int
	for _, key := range primaryKeys {
		if !bytes.Equal(key, primaryKey) {
			cp[j] = key
			j++
		}
	}
	return mi.savePrimaryKeys(secondaryKey, primaryKeys)
}

func (mi *multiIndex) Link(primaryKey, secondaryKey []byte) error {
	mi.lock.Lock()
	defer mi.lock.Unlock()
	primaryKeys, err := mi.getPrimaryKeys(secondaryKey)
	if err != nil {
		return err
	}
	primaryKeys = append(primaryKeys, primaryKey)
	return mi.savePrimaryKeys(secondaryKey, primaryKeys)
}

func (mi *multiIndex) Iterate(handler func(primaryKey, secondaryKey []byte) error) error {
	mi.lock.RLock()
	defer mi.lock.RUnlock()
	return mi.index.Keys(func(secondaryKey []byte) error {
		primaryKeys, err := mi.getPrimaryKeys(secondaryKey)
		if err != nil {
			return err
		}
		for _, primaryKey := range primaryKeys {
			err = handler(primaryKey, secondaryKey)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (mi *multiIndex) getPrimaryKeys(secondaryKey []byte) ([][]byte, error) {
	var primaryKeys [][]byte
	data, err := mi.index.Get(secondaryKey)
	if err == nil {
		err = gob.NewDecoder(bytes.NewReader(data)).Decode(&primaryKeys)
		if err != nil {
			return nil, err
		}
	} else if err != os.ErrNotExist {
		return nil, err
	}
	return primaryKeys, nil
}

func (mi *multiIndex) savePrimaryKeys(secondaryKey []byte, primaryKeys [][]byte) error {
	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(primaryKeys)
	if err != nil {
		return err
	}
	return mi.index.Put(secondaryKey, buf.Bytes())
}
