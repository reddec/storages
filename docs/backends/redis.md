---
backend: "Redis"
package: "std/redistorage"
headline: "Redis hashmap as a storage"
features: ["namespace"]
project_url: "https://github.com/go-redis/redis"
---
{% include backend_head.md page=page %}

Wrapper around Redis hashmap where one storage is one hashmap.

Each namespace is each key supposing that the value is hashmap.

Closing root (parent) storage will close all namespaced storages but not vice-versa.

### URL initialization

Do not forget to import package!

`redis://[[user][:<password>]@]<host>[:port][/<dbnum>][?key=<key>]`

Where:

* `<user>` - optional user name for authorization
* `<password>` - optional password for authorization
* `<host>` - required address of redis database
* `<port>` - optional (default 6379) database port
* `<dbnum>` - optional (default 0) database num
* `<key>` - optional (default `DEFAULT`) name of hashmap to store data 

Example:

* `redis://localhost` - simple, without authorization, on local host
* `redis://user:qwerty@localhost:8899/1?key=data` - with authorization, on localhost with custom port and custom database with hashmap key `data`

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