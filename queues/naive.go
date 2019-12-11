package queues

import (
	"encoding/binary"
	"github.com/pkg/errors"
	"github.com/reddec/storages"
	"os"
	"sync"
)

const (
	latestSequenceKey = "latest"
	oldestSequenceKey = "oldest"
)

// Basic but powerful implementation of queues based on any storage
func NaiveQueue(storage storages.KV) (*naiveQueue, error) {
	oldest, err := loadBinaryKey(storage.Get([]byte(oldestSequenceKey)))
	if err != nil {
		return nil, errors.Wrap(err, "load oldest sequence")
	}
	latest, err := loadBinaryKey(storage.Get([]byte(latestSequenceKey)))
	if err != nil {
		return nil, errors.Wrap(err, "load latest sequence")
	}

	if oldest == 0 {
		oldest = 1
	}

	if oldest > latest && oldest-latest > 1 {
		return nil, errors.Errorf("oldest sequence %v is too much further then latest sequence %v", oldest, latest)
	}

	return &naiveQueue{
		storage:        storage,
		latestSequence: latest,
		oldestSequence: oldest,
	}, nil
}

type naiveQueue struct {
	storage        storages.KV
	latestSequence uint64 // last used sequence ID (0 means unused).
	oldestSequence uint64 // oldest used sequence ID (0 means unused)
	lock           sync.RWMutex
}

func (nq *naiveQueue) Put(data []byte) error {
	nq.lock.Lock()
	defer nq.lock.Unlock()

	num := nq.latestSequence + 1
	sequenceNum := nq.getKey(num)
	// save data
	err := nq.storage.Put(sequenceNum[:], data)
	if err != nil {
		return err
	}
	// save sequence
	err = nq.unsafeWriteLatestSequence(sequenceNum[:])
	if err != nil {
		return err
	}
	// update cached sequence
	nq.latestSequence = num
	return nil
}

func (nq *naiveQueue) Peek() ([]byte, error) {
	nq.lock.RLock()
	defer nq.lock.RUnlock()

	data, _, err := nq.unsafePeek()
	return data, err
}

func (nq *naiveQueue) Get() ([]byte, error) {
	nq.lock.Lock()
	defer nq.lock.Unlock()
	data, _, err := nq.unsafePeek()
	if err != nil {
		return nil, err
	}
	return data, nq.unsafeDiscard()
}

func (nq *naiveQueue) Discard() error {
	nq.lock.Lock()
	defer nq.lock.Unlock()

	return nq.unsafeDiscard()
}

func (nq *naiveQueue) unsafeWriteLatestSequence(currentSequenceID []byte) error {
	return nq.storage.Put([]byte(latestSequenceKey), currentSequenceID)
}

func (nq *naiveQueue) unsafeWriteOldestSequence(currentSequenceID []byte) error {
	return nq.storage.Put([]byte(oldestSequenceKey), currentSequenceID)
}

func (nq *naiveQueue) unsafeDiscard() error {
	if nq.isEmpty() {
		return os.ErrNotExist
	}
	num := nq.oldestSequence
	key := nq.getKey(num)

	err := nq.storage.Del(key[:])
	if err != nil {
		return err
	}

	next := nq.oldestSequence + 1
	sequence := nq.getKey(next)
	err = nq.unsafeWriteOldestSequence(sequence[:])
	if err != nil {
		return err
	}

	nq.oldestSequence = next
	return nil
}

func (nq *naiveQueue) unsafePeek() (data []byte, key [8]byte, err error) {
	if nq.isEmpty() {
		err = os.ErrNotExist
		return
	}

	num := nq.oldestSequence
	key = nq.getKey(num)

	data, err = nq.storage.Get(key[:])
	return
}

func (nq *naiveQueue) isEmpty() bool {
	return nq.latestSequence == 0 || nq.oldestSequence > nq.latestSequence
}

func (nq *naiveQueue) getKey(num uint64) [8]byte {
	var sequenceNum [8]byte
	binary.BigEndian.PutUint64(sequenceNum[:], num)
	return sequenceNum
}

func loadBinaryKey(data []byte, err error) (uint64, error) {
	if err == os.ErrNotExist {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	if len(data) != 8 {
		return 0, errors.Errorf("broken data: required 8 bytes")
	}
	return binary.BigEndian.Uint64(data), nil
}
