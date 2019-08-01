package queues

import (
	"math"
	"os"
	"storages"
	"strconv"
	"sync"
)

type simpleQueue struct {
	storage storages.Accessor
	minID   int64
	nextID  int64
	lock    sync.RWMutex
}

func (sq *simpleQueue) Put(data []byte) (id int64, err error) {
	sq.lock.Lock()
	defer sq.lock.Unlock()
	id = sq.nextID
	key := strconv.FormatInt(id, 10)
	err = sq.storage.Put([]byte(key), data)
	if err != nil {
		return -1, err
	}
	if sq.nextID-sq.minID == 0 {
		// empty queue
		sq.minID = id
	}
	sq.nextID++
	return id, nil
}

func (sq *simpleQueue) Peek() (id int64, data []byte, err error) {
	sq.lock.RLock()
	defer sq.lock.RUnlock()
	id = sq.nextID - 1
	if id < sq.minID {
		// queue is empty
		return -1, nil, os.ErrNotExist
	}
	key := strconv.FormatInt(id, 10)
	data, err = sq.storage.Get([]byte(key))
	return id, data, err
}

func (sq *simpleQueue) Clean(end int64) error {
	sq.lock.Lock()
	defer sq.lock.Unlock()
	if end > sq.nextID {
		end = sq.nextID
	}
	for id := sq.minID; id < end; id++ {
		key := strconv.FormatInt(id, 10)
		err := sq.storage.Del([]byte(key))
		if err != nil {
			return err
		}
		sq.minID++
	}
	return nil
}

func (sq *simpleQueue) Size() int64 {
	sq.lock.RLock()
	defer sq.lock.RUnlock()
	return sq.nextID - sq.minID
}

func (sq *simpleQueue) First() int64 {
	return sq.minID
}

func (sq *simpleQueue) Last() int64 {
	return sq.nextID - 1
}

func (sq *simpleQueue) Iterate(from int64) storages.Iterator {
	sq.lock.RLock()
	defer sq.lock.RUnlock()
	if from == 0 {
		from = sq.minID - 1
	} else {
		from -= 1
	}
	return &simpleIterator{
		current: from,
		storage: sq.storage,
		max:     sq.nextID,
	}
}

// Simple queue based on storage. Scans all keys to find minimum and maximum sequence number
func Simple(storage storages.Storage) (storages.Queue, error) {
	var min int64 = math.MaxInt64
	var max int64 = math.MinInt64
	err := storage.Keys(func(key []byte) error {
		id, err := strconv.ParseInt(string(key), 10, 64)
		if err != nil {
			return nil // skip other keys
		}
		if id > max {
			max = id
		}
		if id < min {
			min = id
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if max == math.MinInt64 {
		max = -1
		min = 0
	}

	return SimpleBound(storage, min, max), nil
}

// Simple queue based on storage with manually defined minimal and maximum sequences
func SimpleBound(storage storages.Accessor, minID, maxID int64) storages.Queue {
	var nextId = maxID + 1
	return &simpleQueue{
		minID:   minID,
		nextID:  nextId,
		storage: storage,
	}
}

type simpleIterator struct {
	max      int64
	current  int64
	finished bool
	value    []byte
	storage  storages.Reader
}

func (si *simpleIterator) ID() int64 {
	return si.current
}

func (si *simpleIterator) Value() []byte {
	return si.value
}

func (si *simpleIterator) Next() bool {
	if si.finished {
		return false
	}
	if si.current >= si.max {
		si.finished = true
		return false
	}
	si.current += 1
	key := strconv.FormatInt(si.current, 10)
	if val, err := si.storage.Get([]byte(key)); err != nil {
		si.finished = true
		return false
	} else {
		si.value = val
	}
	return true
}
