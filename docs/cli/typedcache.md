
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