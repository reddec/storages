# Collection of storages (and wrappers)
[![Documentation](https://img.shields.io/badge/documentation-latest-green)](https://reddec.github.io/storages/)
[![license](https://img.shields.io/github/license/reddec/storages.svg)](https://github.com/reddec/storages)
[![](https://godoc.org/github.com/reddec/storages?status.svg)](http://godoc.org/github.com/reddec/storages)
[![donate](https://img.shields.io/badge/help_by️-donate❤-ff69b4)](http://reddec.net/about/#donate)



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

See [documentation](https://reddec.github.io/storages/) for details

**one more example...**

You can create different storage types just by URL (if you imported required package):

For example, use Redis db as a backend

```go
storage, err := std.Create("redis://localhost")
if err != nil {
    panic(err)
}
defer storage.Close()
```

# License

The wrappers itself licensed under MIT but used libraries may have different license politics.
