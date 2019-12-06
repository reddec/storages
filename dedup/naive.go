package dedup

import (
	"github.com/pkg/errors"
	"github.com/reddec/storages"
	"os"
	"sync"
)

type naive struct {
	lock          sync.RWMutex
	storage       storages.Storage
	keys          int
	maxKeys       int
	cleanupAmount int // maxKeys x2
}

// Naive implementation of deduplicate process: simply keep keys as-is, remove old keys when amount (quantity) increased up to
// maxKeys * cleanFactor till maxKeys count. Relay on Keys() method of storage to detect order of keys.
// Cleaning of old keys initiates in Save() method automatically in a same thread.
func NewNaive(storage storages.Storage, maxKeys int, cleanFactor int) (*naive, error) {
	nv := &naive{
		storage:       storage,
		maxKeys:       maxKeys,
		cleanupAmount: maxKeys * cleanFactor,
	}

	return nv, storage.Keys(func(key []byte) error {
		nv.keys++
		return nil
	})
}

func (nv *naive) IsDuplicated(key []byte) (bool, error) {
	nv.lock.RLock()
	defer nv.lock.RUnlock()

	_, err := nv.storage.Get(key)
	if err == os.ErrNotExist {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

var errEnough = errors.New("enough keys")

func (nv *naive) Save(key []byte) error {
	nv.lock.Lock()
	defer nv.lock.Unlock()
	err := nv.storage.Put(key, []byte(""))
	if err != nil {
		return err
	}
	nv.keys++
	if nv.keys >= nv.cleanupAmount {
		// time to cleanup old keys
		return nv.cleanup()
	}
	return nil
}

func (nv *naive) cleanup() error {
	amountToDelete := nv.keys - nv.maxKeys
	var keysToDelete [][]byte // potential memory overflow, but it's the fastest way
	err := nv.storage.Keys(func(key []byte) error {
		if len(keysToDelete) >= amountToDelete {
			return errEnough
		}
		keysToDelete = append(keysToDelete, key)
		return nil
	})
	if err != nil && err != errEnough {
		return err
	}
	for _, key := range keysToDelete {
		err = nv.storage.Del(key)
		if err != nil {
			return err
		}
		nv.keys--
	}

	return nil
}
