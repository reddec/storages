---
backend: "BBolt"
package: "std/boltdb"
headline: "Single-file, embeddable, pure-Go storage"
features: ["namespace"]
project_url: "https://github.com/etcd-io/bbolt"
---
{% include backend_head.md page=page %}

Generates BoltDB (etc.d fork called bbolt) storage.

Default bucket name is `DEFAULT`. 

Uses buckets as namespaces. Closing root (parent) storage will close all namespaced storages but not vice-versa.

## URL initialization

Do not forget to import package!

`bbolt://<path>`

Where:

* `<path>` - path storage file

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

{% include backend_tail.md page=page %}
