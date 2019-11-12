package storages

import (
	"io"
)

// Sharding pool
type ShardPool interface {
	// Get storage based on key
	Get(key []byte) (Storage, error)
	// Iterate over shards
	Iterate(handler func(storage Storage) error) error
	io.Closer
}

// New sharded storage with defined pool (strategy)
func NewSharded(pool ShardPool) Storage {
	return &shardedStorage{pool: pool}
}

type shardedStorage struct {
	pool ShardPool
}

func (shard *shardedStorage) Put(key []byte, data []byte) error {
	storage, err := shard.pool.Get(key)
	if err != nil {
		return err
	}
	defer storage.Close()
	return storage.Put(key, data)
}

func (shard *shardedStorage) Close() error {
	return shard.pool.Close()
}

func (shard *shardedStorage) Get(key []byte) ([]byte, error) {
	storage, err := shard.pool.Get(key)
	if err != nil {
		return nil, err
	}
	defer storage.Close()
	return storage.Get(key)
}

func (shard *shardedStorage) Del(key []byte) error {
	storage, err := shard.pool.Get(key)
	if err != nil {
		return err
	}
	defer storage.Close()
	return storage.Del(key)
}

func (shard *shardedStorage) Keys(handler func(key []byte) error) error {
	return shard.pool.Iterate(func(storage Storage) error {
		return storage.Keys(handler)
	})
}
