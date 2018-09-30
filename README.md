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