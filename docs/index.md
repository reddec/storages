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

# License

The wrappers itself licensed under MIT but used libraries may have different license politics.

# Code conventions and agreements

* [code and interface style](./convention/coding)

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

* [deduplication](./derived/dedup)
* [queues](./derived/queues)
* [indexes](./derived/indexes)

# CLI 

* [storages](./cli/storages)

## Code-generation

* [typedstorage](./cli/typedstorage)
* [typedcache](./cli/typedcache)


