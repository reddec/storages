package dedup

import (
	"bytes"
	"encoding/binary"
	"github.com/reddec/storages"
	"math/rand"
	"os"
)

// Offloaded deduplication is a wrapper around storage that checks and store keys with random unique iteration id.
// In case if storage does not support Clearable interface, unique random iteration id could be used for distinguish
// keys from different iterations. Good for large data set.
func Offloaded(storage storages.Storage) *offloaded {
	off := &offloaded{
		storage: storage,
	}
	off.reset()
	return off
}

type offloaded struct {
	iterationID []byte
	storage     storages.Storage
}

func (off *offloaded) IsDuplicated(key []byte) (bool, error) {
	offloadedIterationId, err := off.storage.Get(key)
	if err != nil && err != os.ErrNotExist {
		// problem with offload storage
		return false, err
	} else if bytes.Compare(offloadedIterationId, off.iterationID) == 0 {
		// already used key
		return true, nil
	}
	// new key or key not yet recorded for the iteration
	return false, nil
}

func (off *offloaded) Save(key []byte) error {
	return off.storage.Put(key, off.iterationID)
}

func (off *offloaded) Clear() error {
	if cls, ok := off.storage.(storages.Clearable); ok {
		return cls.Clear()
	}
	return nil
}

func (off *offloaded) reset() {
	var iterationID [8]byte
	binary.BigEndian.PutUint64(iterationID[:], rand.Uint64())
	off.iterationID = iterationID[:]
}
