
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
  -k, --key=  Key in storage where configuration defined [$KEY]
  -L, --lock= Optional lock file for inter-process synchronization [$LOCK]

Help Options:
  -h, --help  Show this help message

Available commands:
  config     operations on configuration (aliases: cfg)
  copy       copy keys from storage to destination (aliases: cp, c)
  get        get value by key (aliases: fetch, g)
  list       list keys in storage (aliases: ls)
  queue      access to storage by naive queue interface (aliases: q)
  remove     remove value by key (aliases: delete, del, rm)
  serve      expose storage over REST interface (aliases: rest)
  set        set value for key (aliases: put, s)
  supported  list supported storages backends

```

See `storages <command> --help`

### Queues

```
Usage:
  storages [OPTIONS] queue <command>

Application Options:
  -u, --url=      Storage URL (default: bbolt://data) [$URL]
  -k, --key=      Key in storage where configuration defined [$KEY]
  -L, --lock=     Optional lock file for inter-process synchronization [$LOCK]

Help Options:
  -h, --help      Show this help message

Available commands:
  discard  remove oldest data from queue (like silent get)
  get      get oldest data from queue and remove it (aliases: pop)
  peek     get oldest data from queue but not remove
  put      put data to the queue (aliases: push, append)
  serve    expose queue over REST interface (aliases: rest)
```


# Install

### Binary

Look to releases section

### Debian/Ubuntu

#### Add repository

* supported distribution: trusty, xenial, bionic, buster, wheezy

```bash
echo "deb https://dl.bintray.com/reddec/storages-debian {distribution} main" | sudo tee -a /etc/apt/sources.list
```

**Ubuntu 18.04 (bionic)**

```bash
echo "deb https://dl.bintray.com/reddec/storages-debian bionic main" | sudo tee -a /etc/apt/sources.list
```

**Ubuntu 16.04 (xenial)**

```bash
echo "deb https://dl.bintray.com/reddec/storages-debian xenial main" | sudo tee -a /etc/apt/sources.list
```

#### Update cache

`sudo apt-get update`

### Install

`sudo apt-get install storages`

### Build from source

* requires Go 1.13+

`go get github.com/reddec/storages/cmd/...`
