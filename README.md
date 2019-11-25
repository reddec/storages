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

## Collection


### File storage

import: `github.com/reddec/storages/filestorage`

* `New`, `NewDefault`

Puts each data to separate file. File name generates from hash function (by default SHA256) applied to key. To prevent
generates too much files in one directory, each filename is chopped to 4 slices by 4 characters.

* `NewFlat`

Key is equal to file name. Sub-directories (`/` in key name) are not allowed.

### Level DB

import: `github.com/reddec/storages/leveldbstorage`

Generates LevelDB storage (github.com/syndtr/goleveldb) and stores all item as-is inside DB

### BBolt DB

import: `github.com/reddec/storages/boltdb`

Generates BoltDB (etc.d fork called bbolt) storage

### Memory DB

import: `github.com/reddec/storages/memstorage`

Based on hashmap and RWLock in-memory storage

### NOP

import: `github.com/reddec/storages/memstorage`

No-Operation storage that drops any content and returns not-exists on any request.

Useful for mocking, performance testing or for any other logic that needs discard storage.

### Redis

import: `github.com/reddec/storages/redistorage`

Wrapper around Redis hasmap where one storage is one hashmap.

### S3

import: `github.com/reddec/storages/awsstorage`

Wrapper around official S3 SDK to work with bucket as a map 

# Queues

Wrappers around KV-storage that makes a queues. Idea is to keep minimal and maximum id and use sequence to generate 
next key for KV storage.

import: `github.com/reddec/storages/queues`

## Basic queue

* peek - get last but do not remove
* put - push data to the end of queue
* clean - remove data from the first till specified sequence id. Remove all is: `Clean(queue.Last()+1)`

Constructors:

* `Simple`
* `SimpleBounded`

## Limited queue

Queue that removes old items if no more space available (like circular buffer) on `Put` operation.

Constructors:

* `Limited`
* `SimpleLimited` - shorthand for `Limited(Simple(...))`

# Collection of deduplicate methods

All implementations should follow those interface

```golang
// Deduplicate primitive: check if key is already saved and save key
type Dedup interface {
	// Is key already save?
	IsDuplicated(key []byte) (bool, error)
	// Save key for future checks
	Save(key []byte) (error)
}
```


### Naive


Properties:

* `maxKeys` - maximum keys to store after cleanup
* `cleanFactor` - multiply factor of `maxKeys` that triggers cleanup process

Naive implementation of deduplicate process: simply keep keys as-is, remove old keys when amount (quantity) increased up to
`maxKeys * cleanFactor` till `maxKeys count`.

Relay on storages to detect order of keys.

Cleaning of old keys initiates in `Save()`` method automatically in a same thread.

## CLI access

```go get -v github.com/reddec/storages/cmd/storages```

There is a simple command line wrapper around all currently supported database: `storage`. It provides get, put, del and
list operations over file, leveldb and redis storage.

Important: empty value implies stream.

Usage:

```
Usage:
  storages [OPTIONS] [Command] [key] [Value]

Application Options:
  -t, --db=[file|leveldb|redis|s3] DB mode (default: file) [$DB]
  -s, --stream                     Use STDIN as source of value [$STREAM]
  -0, --null                       Use zero byte as terminator for list instead of new line [$NULL]

File storage params:
      --file.location=             Root dir to store data (default: ./db) [$FILE_LOCATION]

LevelDB storage params:
      --leveldb.location=          Root dir to store data (default: ./db) [$LEVELDB_LOCATION]

Redis storage params:
      --redis.url=                 Redis URL (default: redis://localhost) [$REDIS_URL]
      --redis.namespace=           Hashmap name (default: db) [$REDIS_NAMESPACE]

S3 storage:
      --s3.bucket=                 S3 AWS bucket [$S3_BUCKET]
      --s3.endpoint=               Override AWS endpoint for AWS-capable services [$S3_ENDPOINT]
      --s3.force-path-style        Force the request to use path-style addressing [$S3_FORCE_PATH_STYLE]

Help Options:
  -h, --help                       Show this help message

Arguments:
  Command:                         what to do (put, list, get, del)
  key:                             key name
  Value:                           Value to put if stream flag is not enabled

```

## CLI generators

Supports code-generation

### typedstorage


Typed wrapper around any storage with JSON encoding by default

Usage:

    Usage of typedstorage:
      -out string
            Output file (default: <type name>_storage.go)
      -package string
            Output package (default: same as in input file)
      -type string
            Type name to wrap


Embedded usage example:

```go


type Sample struct {
    // ...
}

//go:generate typedstorage -type Sample

```

will produce (methods body omitted, see sample dir for details)


```go
// Typed storage for Sample
type SampleStorage struct {
	cold storages.Storage // persist storage
}

// Creates new storage for Sample
func NewSampleStorage(cold storages.Storage) *SampleStorage {}

// Put single Sample encoded in JSON into storage
func (cs *SampleStorage) Put(key string, item *Sample) error {}

// Get single Sample from storage and decode data as JSON
func (cs *SampleStorage) Get(key string) (*Sample, error) {}

// Del key from hot and cold storage
func (cs *SampleStorage) Del(key string) error {}

// Keys copied slice that cached in hot storage
func (cs *SampleStorage) Keys() ([]string, error) {}

```


### typedcache


Typed wrapper around any storage with JSON encoding by default with in-memory cache.

    Usage of typedstorage:
      -codec string
            Encoder/Decoder for the type: json, msgp (default "json")
      -out string
            Output file (default: <type name>_storage.go)
      -package string
            Output package (default: same as in input file)
      -prefix string
            Custom key prefix
      -type string
            Type name to wrap


Note: `msgp` codec requires msgp binding

Embedded usage example:

```go

type Sample struct {
    // ...
}

//go:generate typedcache -type Sample

```

will produce (methods body omitted, see sample dir for details)


```go

// Two level storage for Sample
type CachedSampleStorage struct {
	cold storages.Storage   // persist storage
	hot  map[string]*Sample // cache storage
	lock sync.RWMutex
}

// Creates new storage for Sample with custom cache
func NewCachedSampleStorage(cold storages.Storage) *CachedSampleStorage {}

// Put single Sample encoded in JSON into cold and hot storage
func (cs *CachedSampleStorage) Put(key string, item *Sample) error {}

/*
Get single Sample from hot storage and decode data as JSON.
If key is not in hot storage, the cold storage is used and obtained data is put to the hot storage for future cache
*/
func (cs *CachedSampleStorage) Get(key string) (*Sample, error) {}

// Fetch all data from cold storage to the hot storage (warm cache)
func (cs *CachedSampleStorage) Fetch() error {}

// Keys copied slice that cached in hot storage
func (cs *CachedSampleStorage) Keys() []string {}

// Del key from hot and cold storage
func (cs *CachedSampleStorage) Del(key string) error {}

func (cs *CachedSampleStorage) getMissed(key string) (*Sample, error) {}

```
