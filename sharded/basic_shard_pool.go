package sharded

import (
	"github.com/reddec/storages"
	"hash/crc32"
	"sync"
)

// Hash function of key to determinate shard id. Non-scaled
type HashShardFunc func([]byte) uint32

// Factory function to create storage by id (id scaled in range of defined shard number)
type StorageFactoryFunc func(shardID uint32) (storages.Storage, error)

// New shard pool based on default hash distribution IEEE CRC-32. See HashedCustom for details.
func NewHashed(shards uint32, factory StorageFactoryFunc) *hashShard {
	return NewHashedCustom(shards, factory, crc32.ChecksumIEEE)
}

// New  shard pool based on custom hash distribution.
// Shards will be open dynamically on-demand and kept open till pool closed.
// Thread safe.
func NewHashedCustom(shards uint32, factory StorageFactoryFunc, hashFunc HashShardFunc) *hashShard {
	if shards == 0 {
		panic("zero shards defined")
	}
	cache := make([]*shardItem, shards)
	for i := range cache {
		cache[i] = new(shardItem)
	}
	return &hashShard{
		shards:   cache,
		hashFunc: hashFunc,
		factory:  factory,
	}
}

type shardItem struct {
	storage storages.Storage
	lock    sync.Mutex
}

type hashShard struct {
	shards   []*shardItem
	hashFunc HashShardFunc
	factory  StorageFactoryFunc
}

func (hs *hashShard) Get(key []byte) (storages.Storage, error) {
	shardID := hs.hashFunc(key) % uint32(len(hs.shards))
	return hs.getOrCreate(shardID)
}

func (hs *hashShard) Iterate(handler func(storage storages.Storage) error) error {
	for i := 0; i < len(hs.shards); i++ {
		storage, err := hs.getOrCreate(uint32(i))
		if err != nil {
			return err
		}
		err = handler(storage)
		if err != nil {
			return err
		}
	}
	return nil
}

func (hs *hashShard) Close() error {
	for _, shard := range hs.shards {
		shard.lock.Lock()
		if shard.storage != nil {
			shard.storage.Close()
		}
		shard.lock.Unlock()
	}
	return nil
}

func (hs *hashShard) getOrCreate(shardID uint32) (storages.Storage, error) {
	shard := hs.shards[shardID]
	if shard.storage != nil {
		return &noClose{shard.storage}, nil
	}
	shard.lock.Lock()
	defer shard.lock.Unlock()
	if shard.storage != nil {
		return &noClose{shard.storage}, nil
	}
	storage, err := hs.factory(shardID)
	if err != nil {
		return nil, err
	}
	shard.storage = storage
	return &noClose{shard.storage}, nil
}

type noClose struct {
	storage storages.Storage
}

func (n *noClose) Close() error                              { return nil }
func (n *noClose) Put(key []byte, data []byte) error         { return n.storage.Put(key, data) }
func (n *noClose) Get(key []byte) ([]byte, error)            { return n.storage.Get(key) }
func (n *noClose) Del(key []byte) error                      { return n.storage.Del(key) }
func (n *noClose) Keys(handler func(key []byte) error) error { return n.storage.Keys(handler) }
