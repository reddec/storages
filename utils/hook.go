package stor_utils

import (
	"github.com/reddec/storages"
)

// Wrap storage with additional hook that will be called before close
func WithCloseHook(storage storages.Storage, hook func()) *withCloser {
	return &withCloser{
		stor: storage,
		hook: hook,
	}
}

type withCloser struct {
	stor storages.Storage
	hook func()
}

func (f *withCloser) Close() error {
	f.hook()
	return f.stor.Close()
}

func (f *withCloser) Put(key []byte, data []byte) error {
	return f.stor.Put(key, data)
}

func (f *withCloser) Get(key []byte) ([]byte, error) {
	return f.stor.Get(key)
}

func (f *withCloser) Del(key []byte) error {
	return f.stor.Del(key)
}

func (f *withCloser) Keys(handler func(key []byte) error) error {
	return f.stor.Keys(handler)
}
