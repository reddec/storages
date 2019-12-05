---
backend: "REST"
headline: "REST-like storage with server handler"
package: "std/rest"
features: []
project_url: ""
---
{% include backend_head.md page=page %}

REST-like storage. Expects keys in standard base-64 encoding. Server interface should be like described:

| Method   | Path       | Success status | Description |
|----------|------------|----------------|-------------
| `GET`    | `/`        | 200            | Array of all keys. New line - new key.
| `GET`    | `/:key`    | 200            | Content of key as-is without encoding
| `POST`   | `/:key`    | 204            | Update or insert value for key
| `DELETE` | `/:key`    | 204            | Remove key. Removing non-existent should be not an error

Returned status as not the same as expected in success column means error. **Special case** for `GET /:key` when there is
no key in a storage, but operation successful:  `404` MUST be returned.

**Expose storage**

You may expose any storage that follow `Storage` interface by simple wrapper: `NewServer(storage)`

## Usage

**Example client**

```go
storage := rest.New("http://example.com:8080/path/to/storage")
// Get, Keys and so on
```

**Example server**

* `myStorage` is somewhere defined in your code storage instance 
* do not forget to add trailing slash (`/`) in export path, or nested routes will not be routed

```go
const exportPath = "/path/to/storage/"
handler := rest.NewServer(myStorage)

http.Handle(exportPath, http.StripPrefix(exportPath, handler))
// ...
http.ListenAndServe(":8080", nil) 
```

{% include backend_tail.md page=page %}