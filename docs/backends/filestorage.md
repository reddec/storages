---
backend: "Filesystem"
package: "filestorage"
headline: "Local file-system storage"
features: ["namespace"]
project_url: ""
---

{% include backend_head.md page=page %}

### Encoded 

Constructors: `New`, `NewDefault`

Puts each data to separate file. File name generates from hash function (by default SHA256) applied to key. To prevent
generates too much files in one directory, each filename is chopped to 4 slices by 4 characters.


### Flat

Constructor: `NewFlat`

Key is equal to file name. Sub-directories (`/` in key name) are not allowed.

Namespace are share key space with regular values.

## Usage

**Flat**

```go
stor, err = filestorage.NewFlat("path/to/directory")
if err != nil {
    panic(err)
}
// Close() not required but implemented as NOP
```


**Encoded**

```go
stor, err = filestorage.NewDefault("path/to/directory")
if err != nil {
    panic(err)
}
// Close() not required but implemented as NOP
```


{% include backend_tail.md page=page %}
