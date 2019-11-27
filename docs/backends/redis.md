# Redis

[![API docs](https://godoc.org/github.com/reddec/storages/redistorage?status.svg)](http://godoc.org/github.com/reddec/storages/redistorage)

* **import:** `github.com/reddec/storages/redistorage`
* [Redis project](https://github.com/go-redis/redis) 

Wrapper around Redis hashmap where one storage is one hashmap.

Each namespace is each key supposing that the value is hashmap.

Closing root (parent) storage will close all namespaced storages but not vice-versa.

## Features

{% include feature_namespace.md %}

## Examples

**With new connection**

```go
storage, err := New("my-space", "redis://redishost")
if err != nil {
   panic(err)    
}
defer storage.Close()
```

or with helper that will panic on error

```go
storage := MustNew("my-space", "redis://redishost")
defer storage.Close()
```

**Wrap connection**

```go
storage := NewClient("my-space", redisConnection)
defer storage.Close()
```

