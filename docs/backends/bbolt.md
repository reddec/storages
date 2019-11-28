---
backend: "BBolt"
headline: "Single-file, embeddable, pure-Go storage"
features: ["namespace"]
---
# BBolt DB

[![API docs](https://godoc.org/github.com/reddec/storages/boltdb?status.svg)](http://godoc.org/github.com/reddec/storages/boltdb)

* **import:** `github.com/reddec/storages/boltdb`
* [BBolt project](https://github.com/etcd-io/bbolt)

Generates BoltDB (etc.d fork called bbolt) storage.

Default bucket name is `DEFAULT`. 

Uses buckets as namespaces. Closing root (parent) storage will close all namespaced storages but not vice-versa.



## Features

{% include feature_namespace.md %}

## Usage

**With default bucket**

```go
storage, err := boltdb.NewDefault("path/to/file")
if err != nil {
    panic(err)
}
defer storage.Close()
```

**With custom bucket**

```go
storage, err := boltdb.New("path/to/file", []byte("custom bucket name"))
if err != nil {
    panic(err)
}
defer storage.Close()
```