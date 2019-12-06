# Sharding

Sharding is a process to distribute data through multiple storages. 

                  storage1 (data from 0..x)
                🡕
    data(0...N) 🡒 storage2 (data from 0..y)
                🡖 
                  storageM (data from 0..z)
    

Sharded storage mimics to usual [Storage](https://godoc.org/github.com/reddec/storages#Storage) interface so 
target consumers should work with it as usual.

Constructor for sharding storage is [Sharded(pool)](https://godoc.org/github.com/reddec/storages#Sharded).

Distribution logic is pluggable by [ShardPool](https://godoc.org/github.com/reddec/storages#ShardPool) that responsible
for storage allocation and caching (if needed).

Usage is quite straightforward:

```go
shardedStorage := Sharded(pool)
defer shardedStorage.Close()
// then as usual storage
```

There are some pre-defined pool implementations

## Hashed

**import:** `github.com/reddec/storages/sharded`

Hashed pool bases on some integer hash calculated on key and then scaled to number of shard.

    KEY - key bytes
    N - number of shards
    
    H(key) -> :u32  // hash functin, where result in range [0..2^32-1)
    S(key) -> :u32  // shard id, where id in range [0..N)
    
    S(key) = H(key) % N 

By default `sharded.New` is using `CRC32-IEEE` as hash function. It could be enough for testing for several cases
it could be not enough.

It's possible to use any custom hash function by using `sharded.NewCustom`.

Number of shards should be constant. If you need to change number of shards you have to make a new
sharded storage with new number of shard and copy data from old storage.  

### Usage

**In-memory**

Shards num: 3

```go
pool := sharded.NewHashed(3, func(shardID uint32) (storage storages.Storage, e error) {
    return memstorage.New(), nil
})

shardedStorage := storages.Sharded(pool)
defer shardedStorage.Close()

```


**BBolts DB as storage**

Shards num: 3

```go
pool := sharded.NewHashed(3, func(shardID uint32) (storage storages.Storage, e error) {
    return boltdb.NewDefault(fmt.Sprtin("shard-", shardID))
})

shardedStorage := storages.Sharded(pool)
defer shardedStorage.Close()

```

**Different storages**

* Shard 1 - bbolt DB
* Shard 2 - level DB
* Shards num: 2

```go
pool := sharded.NewHashed(2, func(shardID uint32) (storage storages.Storage, e error) {
    if shardID == 0 {
        return boltdb.NewDefault(fmt.Sprtin("shard-", shardID))    
    }
    return leveldbstorage.New(fmt.Sprtin("shard-", shardID))
})

shardedStorage := storages.Sharded(pool)
defer shardedStorage.Close()

```