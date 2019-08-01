package memstorage

import (
	"os"
	"storages"
)

type nopStorage struct{}

func (np *nopStorage) Put(key []byte, data []byte) error         { return nil }
func (np *nopStorage) Get(key []byte) ([]byte, error)            { return nil, os.ErrNotExist }
func (np *nopStorage) Del(key []byte) error                      { return nil }
func (np *nopStorage) Keys(handler func(key []byte) error) error { return nil }
func (np *nopStorage) Close() error                              { return nil }
func (np *nopStorage) BatchWriter() storages.Writer              { return NewNOP() }

// New No-Operation storage that drops any content and returns not-exists on any request.
// Useful for mocking, performance testing or for dropping several keys.
func NewNOP() storages.BatchedStorage {
	return &nopStorage{}
}
