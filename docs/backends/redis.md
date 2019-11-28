---
backend: "Redis"
package: "redistorage"
headline: "Redis hashmap as a storage"
features: ["namespace"]
project_url: "https://github.com/go-redis/redis"
---
{% include backend_head.md page=page %}

Wrapper around Redis hashmap where one storage is one hashmap.

Each namespace is each key supposing that the value is hashmap.

Closing root (parent) storage will close all namespaced storages but not vice-versa.


## Usage

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

{% include backend_tail.md page=page %}