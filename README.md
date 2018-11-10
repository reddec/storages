# Collection of storages

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

## Collection


### File storage

import: `github.com/reddec/storages/filestorage`

Puts each data to separate file. File name generates from hash function (by default SHA256) applied to key. To prevent
generates too much files in one directory, each filename is chopped to 4 slices by 4 characters.


### Level DB

import: `github.com/reddec/storages/leveldbstorage`

Generates LevelDB storage (github.com/syndtr/goleveldb) and stores all item as-is inside DB

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
