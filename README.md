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

See [documentation](https://reddec.github.io/storages/) for details

**one more example...**

You can create different storage types just by URL (if you imported required package):

For example, use Redis db as a backend

```go
storage, err := std.Create("redis://localhost")
if err != nil {
    panic(err)
}
defer storage.Close()
```

# CLI tools

### Binary

Look to releases section

### Debian/Ubuntu

[![Debian version](https://api.bintray.com/packages/reddec/storages-debian/storages/images/download.svg)](https://bintray.com/reddec/storages-debian/storages/_latestVersion)

Add public Bintray key

```bash
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 379CE192D401AB61
```

#### Add repository

* supported distribution: trusty, xenial, bionic, buster, wheezy

```bash
echo "deb https://dl.bintray.com/reddec/debian {distribution} main" | sudo tee -a /etc/apt/sources.list
```

**Ubuntu 18.04 (bionic)**

```bash
echo "deb https://dl.bintray.com/reddec/debian bionic main" | sudo tee -a /etc/apt/sources.list
```

**Ubuntu 16.04 (xenial)**

```bash
echo "deb https://dl.bintray.com/reddec/debian xenial main" | sudo tee -a /etc/apt/sources.list
```

#### Update cache

`sudo apt-get update`

### Install

`sudo apt-get install storages`

### Build from source

* requires Go 1.13+

`go get github.com/reddec/storages/cmd/...`

# License

The wrappers itself licensed under MIT but used libraries may have different license politics.
