package storages

import (
	"github.com/pkg/errors"
	"os"
	"strings"
	"sync"
)

// Distributed writer strategy
type DWriter func(key, data []byte, storages []Storage) error

// Distributed reader strategy
type DReader func(key []byte, storages []Storage) ([]byte, error)

// Redundant storage that writes values to all back storages and read from first successful.
func RedundantAll(keysDeduplication Dedup, back ...Storage) *redundant {
	return Redundant(AtLeast(len(back)), First(), keysDeduplication, back...)
}

// Redundant storage with custom strategy for writing and reading backed by several storage
func Redundant(writer DWriter, reader DReader, keysDeduplication Dedup, back ...Storage) *redundant {
	return &redundant{
		backed:            back,
		writer:            writer,
		reader:            reader,
		keysDeduplication: keysDeduplication,
	}
}

type redundant struct {
	backed            []Storage // storages for data
	keysDeduplication Dedup     // used for deduplication during iteration
	writer            DWriter
	reader            DReader
	iterationLock     sync.Mutex
}

func (dt *redundant) Put(key []byte, data []byte) error {
	return dt.writer(key, data, dt.backed)
}

func (dt *redundant) Get(key []byte) ([]byte, error) {
	return dt.reader(key, dt.backed)
}

func (dt *redundant) Close() error {
	var list []error
	for _, stor := range dt.backed {
		list = append(list, stor.Close())
	}
	return allErr(list...)
}

func (dt *redundant) Del(key []byte) error {
	var list []error
	for _, stor := range dt.backed {
		err := stor.Del(key)
		if err != nil {
			list = append(list, err)
		}
	}
	return allErr(list...)
}

func (dt *redundant) Keys(handler func(key []byte) error) error {
	dt.iterationLock.Lock()
	defer dt.iterationLock.Unlock()
	var list []error
	for _, stor := range dt.backed {
		err := stor.Keys(func(key []byte) error {
			isExists, err := dt.keysDeduplication.IsDuplicated(key)
			if err != nil {
				return err
			}
			if isExists {
				return nil
			}
			err = dt.keysDeduplication.Save(key)
			if err != nil {
				return err
			}
			return handler(key)
		})
		if err != nil {
			list = append(list, err)
		}
	}
	// clean prev offload if possible
	if clearable, ok := dt.keysDeduplication.(Clearable); ok {
		_ = clearable.Clear()
	}
	return allErr(list...)
}

func allErr(list ...error) error {
	var ans []string
	for _, err := range list {
		if err != nil {
			ans = append(ans, err.Error())
		}
	}
	if len(ans) == 0 {
		return nil
	}
	return errors.New(strings.Join(ans, "; "))
}

// strategies for read/write/dedup

// At least minWrite amount of written operations should be complete for success result
func AtLeast(minWrite int) DWriter {
	return func(key, data []byte, storages []Storage) error {
		var wrote int
		var list []error
		for _, stor := range storages {
			err := stor.Put(key, data)
			if err != nil {
				list = append(list, err)
			} else {
				wrote++
			}
		}
		if wrote < minWrite {
			return allErr(list...)
		}
		return nil
	}
}

// Shorthand for AtLeast(1) - requires at least one successful write operation
func Any() DWriter { return AtLeast(1) }

// First non-empty value for key will be used as result
func First() DReader {
	return func(key []byte, storages []Storage) ([]byte, error) {
		for _, stor := range storages {
			data, err := stor.Get(key)
			if err == nil {
				return data, nil
			}
		}
		return nil, os.ErrNotExist
	}
}
