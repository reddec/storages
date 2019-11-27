
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

