# Collection of storages (and wrappers)

[![license](https://img.shields.io/github/license/reddec/storages.svg)](https://github.com/reddec/storages)
[![](https://godoc.org/github.com/reddec/storages?status.svg)](http://godoc.org/github.com/reddec/storages)
[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=4UKBSN5HVB3Y8&source=url)

Donating always welcome

* ETH: `0xA4eD4fB5805a023816C9B55C52Ae056898b6BdBC`
* BTC: `bc1qlj4v32rg8w0sgmtk8634uc36evj6jn3d5drnqy`


Different implementations of storages with same abstract interface:


```go
// Thread-safe storage for key-value
type Storage interface {
	// Put single item to storage. If already exists - override
	Put(key []byte, data []byte) error
	// Get item from storage. If not exists - os.ErrNotExist (implementation independent)
	Get(key []byte) ([]byte, error)
	// Delete key and value
	Del(key []byte) error
	// Iterate over all keys. Modification during iteration may cause undefined behaviour (mostly - dead-lock)
	Keys(handler func(key []byte) error) error
    // Close storage if needs
    io.Closer
}
```

# License

The wrappers itself licensed under MIT but used libraries may have different license politics.

# Coding style

* V0 - return abstract interfaces
* V1 - follow 'accept interfaces, return structs'

Since V1 all implementations should return non-exported reference to structure (see `boltdb` wrapper as an example). Standard wrappers will be replace as sooner as possible, 
however it should not affect code that already using current library.

# Backends

* [BBolt](./backends/bbolt)
* [LevelDB](./backends/leveldb)
* [Memory](./backends/memory)
* [Mock](./backends/mock)
* [Redis](./backends/redis)
* [Files](./backends/filestorage.md)
* [S3](./backends/s3.md)

# Derived 

* [deduplication](./derived/dedup)
* [queues](./derived/queues)

## CLI 

* [storages](./cli/storages)

Code-generation

* [typedstorage](./cli/typedstorage)
* [typedcache](./cli/typedcache)



