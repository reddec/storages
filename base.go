package storages

import "io"

// Key-value writer
type Writer interface {
	// Put single item to storage. If already exists - override
	Put(key []byte, data []byte) error
	// Close storage if needs
	io.Closer
}

// Thread-safe storage for key-value
type KV interface {
	Writer
	// Get item from storage. If not exists - os.ErrNotExist (implementation independent)
	Get(key []byte) ([]byte, error)
	// Delete key and value
	Del(key []byte) error
}

// Extension for KV storage with iterator over keys
type Storage interface {
	KV
	// Iterate over all keys. Modification during iteration may cause undefined behaviour (mostly - dead-lock)
	Keys(handler func(key []byte) error) error
}

// Atomic (batch) writer. Batch storage should be used only in one thread
type BatchedStorage interface {
	Storage
	BatchWriter() Writer
}

// Extract all keys from storage as-is
func AllKeys(storage Storage) ([][]byte, error) {
	if storage == nil {
		return nil, nil
	}
	var ans [][]byte
	err := storage.Keys(func(key []byte) error {
		ans = append(ans, key)
		return nil
	})
	return ans, err
}

// Extract all keys from storage and convert it to string
func AllKeysString(storage Storage) ([]string, error) {
	if storage == nil {
		return nil, nil
	}
	var ans []string
	err := storage.Keys(func(key []byte) error {
		ans = append(ans, string(key))
		return nil
	})
	return ans, err
}
