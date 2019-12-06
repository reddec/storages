package storages

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
	"math/rand"
	"os"
	"strings"
)

func RedundantAll(keysOffload Storage, back ...Storage) *redundant {
	return Redundant(len(back), keysOffload, back...)
}

func Redundant(minWrite int, keysOffload Storage, back ...Storage) *redundant {
	return &redundant{
		backed:      back,
		minWrite:    minWrite,
		keysOffload: keysOffload,
	}
}

type redundant struct {
	backed      []Storage // storages for data
	minWrite    int       // minimal amount of written storage for success
	keysOffload Storage   // storage used for deduplication during iteration
}

func (dt *redundant) Put(key []byte, data []byte) error {
	var wrote int
	var list []error
	for _, stor := range dt.backed {
		err := stor.Put(key, data)
		if err != nil {
			list = append(list, err)
		} else {
			wrote++
		}
	}
	if wrote < dt.minWrite {
		return allErr(list...)
	}
	return nil
}

func (dt *redundant) Close() error {
	var list []error
	for _, stor := range dt.backed {
		list = append(list, stor.Close())
	}
	return allErr(list...)
}

func (dt *redundant) Get(key []byte) ([]byte, error) {
	var list []error
	for _, stor := range dt.backed {
		data, err := stor.Get(key)
		if err != nil {
			list = append(list, err)
		} else {
			return data, nil
		}
	}
	return nil, allErr(list...)
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
	var list []error
	// generate unique id for iteration for keys offloading
	var iterationID [8]byte
	binary.BigEndian.PutUint64(iterationID[:], rand.Uint64())
	// clean prev offload if possible
	if clearable, ok := dt.keysOffload.(Clearable); ok {
		err := clearable.Clear()
		if err != nil {
			return err
		}
	}
	for _, stor := range dt.backed {
		err := stor.Keys(func(key []byte) error {
			offloadedIterationId, err := dt.keysOffload.Get(key)
			if err != nil && err != os.ErrNotExist {
				// problem with offload storage
				return err
			} else if bytes.Compare(offloadedIterationId, iterationID[:]) == 0 {
				// already used key
				return nil
			}
			// new key or key not yet recorded for the iteration
			err = dt.keysOffload.Put(key, iterationID[:])
			if err != nil {
				return err
			}
			return handler(key)
		})
		if err != nil {
			list = append(list, err)
		}
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
