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


## Example usage

**static backend**

```go
func run() {
	storage :=  memstorage.New()
	storage.Put([]byte("alice"), []byte("data"))
	restored, _ := storage.Get([]byte("alice"))
	fmt.Println("Restored:", string(restored))
}

```

Due to standardized interface you may use

**dynamic backend selection**

(error handling omitted)

```go
func run(backendName string) {
	storage, _ := create(backendName)
	defer storage.Close()
	storage.Put([]byte("alice"), []byte("data"))
	restored, _ := storage.Get([]byte("alice"))
	fmt.Println("Restored:", string(restored))
}

func create(backend string) (storages.Storage, error) {
	switch backend {
	case "file":
		return filestorage.NewDefault("data"), nil
	case "bolt", "bbolt", "bboltdb":
		return boltdb.NewDefault("data")
	case "memory":
		return memstorage.New(), nil
	default:
		return memstorage.NewNOP(), nil
	}
}

```

# Backends and features

Table of all supported backends and their features.

{%- assign features = [] %}
{%- for page in site.pages %}
{%- if page.dir contains "/backends/" %}
{%- assign features = features | concat: page.features %}
{%- endif %}
{%- endfor %}
{%- assign features = features | sort | uniq %}

|  Backend  | Description   | {{features | join: " | "}}   |
|-----------|---------------|{%for feature in features %}:------------:|{%endfor%}
{%- for page in site.pages %}
{%- if page.dir contains "/backends/" %}
|  [{{page.backend}}]({{page.url | relative_url}})  |  {{page.headline}} {%for feature in features %} | {% if page.features contains feature %} ✔ {%endif%} {%endfor%}  |
{%- endif %}
{%- endfor %}

# Derived 

* [deduplication](./derived/dedup) - deduplication by key
* [queues](./derived/queues) - make queue with any storage as backend
* [sharding](./derived/sharding) - make storage that will distribute values to the different shard 
* [indexes](./derived/indexes) - secondary unique and non-unique indexes

# CLI 

* [storages](./cli/storages)

## Code-generation

* [typedstorage](./cli/typedstorage) - make compile time type-safe wrapper around storage
* [typedcache](./cli/typedcache) - make type-safe wrapper storage with hot and cold layer

# Code conventions and agreements

* [code and interface style](./convention/coding)


# License

The wrappers itself licensed under MIT but used libraries may have different license politics.

