---
backend: "LevelDB"
headline: "Multi-files, embeddable, pure-Go storage"
---
# Level DB

[![API docs](https://godoc.org/github.com/reddec/storages/leveldbstorage?status.svg)](http://godoc.org/github.com/reddec/storages/leveldbstorage)

* **import:** `github.com/reddec/storages/leveldbstorage`
* [LevelDB project](https://github.com/syndtr/goleveldb) 

Multi-files, embeddable, pure-Go storage. Uses levelDB storage as backend. Supports native batching.

## Usage

```go
stor, err = leveldbstorage.New("path/to/dbdir")
if err != nil {
    panic(err)
}
defer stor.Close()
```

## Features

{% include feature_batch_writer.md %}