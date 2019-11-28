---
backend: "LevelDB"
package: "leveldbstorage"
headline: "Multi-files, embeddable, pure-Go storage"
features: ["batch_writer"]
project_url: "https://github.com/syndtr/goleveldb"
---
{% include backend_head.md page=page %}

Multi-files, embeddable, pure-Go storage. Uses levelDB storage as backend. Supports native batching.

## Usage

```go
stor, err = leveldbstorage.New("path/to/dbdir")
if err != nil {
    panic(err)
}
defer stor.Close()
```

{% include backend_tail.md page=page %}