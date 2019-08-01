package memstorage

import (
	"os"
	"storages"
	"sync"
)

type memoryMap struct {
	db   map[string][]byte
	lock sync.RWMutex
}

func (bdp *memoryMap) Put(key []byte, value []byte) error {
	bdp.lock.Lock()
	defer bdp.lock.Unlock()
	if bdp.db == nil {
		bdp.db = make(map[string][]byte)
	}
	k := string(key)
	cp := make([]byte, len(value))
	copy(cp, value)
	bdp.db[k] = cp
	return nil
}

func (bdp *memoryMap) Get(key []byte) ([]byte, error) {
	bdp.lock.RLock()
	defer bdp.lock.RUnlock()
	k := string(key)
	value, ok := bdp.db[k]
	if !ok {
		return nil, os.ErrNotExist
	}
	cp := make([]byte, len(value))
	copy(cp, value)
	return value, nil
}

func (bdp *memoryMap) Del(key []byte) error {
	bdp.lock.Lock()
	defer bdp.lock.Unlock()
	k := string(key)
	delete(bdp.db, k)
	return nil
}

func (bdp *memoryMap) Keys(handler func(key []byte) error) error {
	bdp.lock.RLock()
	defer bdp.lock.RUnlock()
	for k := range bdp.db {
		err := handler([]byte(k))
		if err != nil {
			return err
		}
	}
	return nil
}

func (ds *memoryMap) Close() error { return nil } // NOP

type memBatch struct {
	data map[string][]byte
	mm   *memoryMap
}

func (mb *memBatch) Put(key []byte, value []byte) error {
	cp := make([]byte, len(value))
	copy(cp, value)
	mb.data[string(key)] = cp
	return nil
}

func (mb *memBatch) Close() error {
	mb.mm.lock.Lock()
	defer mb.mm.lock.Unlock()
	if mb.mm.db == nil {
		mb.mm.db = make(map[string][]byte)
	}
	for k, v := range mb.data {
		mb.mm.db[k] = v
	}
	mb.data = nil
	return nil
}

func (ds *memoryMap) BatchWriter() storages.Writer {
	return &memBatch{data: make(map[string][]byte), mm: ds}
}

// New in-memory storage, based on Go concurrent map. For each Add and Get new copy of key and data will be made.
func New() storages.BatchedStorage {
	return &memoryMap{}
}
