
## CLI access

```go get -v github.com/reddec/storages/cmd/storages```

There is a simple command line wrapper around all currently supported database: `storage`. It provides get, put, del and
list operations over file, leveldb and redis storage.

Important: empty value implies stream.

Usage:

```
Usage:
  storages [OPTIONS] <command>

Application Options:
  -u, --url=  Storage URL (default: bbolt://data) [$URL]

Help Options:
  -h, --help  Show this help message

Available commands:
  get        get value by key (aliases: fetch, g)
  list       list keys in storage (aliases: ls)
  remove     remove value by key (aliases: delete, del, rm)
  set        set value for key (aliases: put, s)
  supported  list supported storages backends

```

See `storages <command> --help`