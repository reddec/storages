package sharded

import (
	"errors"
	"github.com/reddec/storages"
)

func ExampleNewHashed() {
	pool := NewHashed(3, func(shardID uint32) (storages.Storage, error) {
		return nil, errors.New("TODO: any storage implementation")
	})

	shardedStorage := storages.Sharded(pool)
	defer shardedStorage.Close()
	// do something
}
