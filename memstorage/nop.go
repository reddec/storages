package memstorage

import (
	"github.com/reddec/storages"
	"os"
)

type nopStorage struct{}

func (np *nopStorage) Put(key []byte, data []byte) error         { return nil }
func (np *nopStorage) Get(key []byte) ([]byte, error)            { return nil, os.ErrNotExist }
func (np *nopStorage) Del(key []byte) error                      { return nil }
func (np *nopStorage) Keys(handler func(key []byte) error) error { return nil }
func (np *nopStorage) Close() error                              { return nil }

// New No-Operation storage that drops any content and returns not-exists on any request.
// Useful for mocking, performance testing or for dropping several keys.
func NewNOP() storages.Storage {
	return &nopStorage{}
}
