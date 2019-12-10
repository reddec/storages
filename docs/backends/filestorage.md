---
backend: "Filesystem"
package: "std/filestorage"
headline: "Local file-system storage"
features: ["namespace"]
project_url: ""
---

{% include backend_head.md page=page %}

### Encoded 

Constructors: `New`, `NewDefault`

Puts each data to separate file. File name generates from hash function (by default SHA256) applied to key. To prevent
generates too much files in one directory, each filename is chopped to 4 slices by 4 characters.

#### URL initialization

Do not forget to import package!

`file://<path>`

Where:

* `<path>` - path to root directory


### Flat

Constructor: `NewFlat`

Key is equal to file name. Sub-directories (`/` in key name) are not allowed.

Namespace are share key space with regular values.

#### URL initialization

Do not forget to import package!

`file+flat://<path>`

Where:

* `<path>` - path to root directory

### JSON

Constructor: `NewJSONFile`

Stores everything in one file. File overwrites atomically. 

With JSON encoding all data presented as dictionary in `data` field and namespaces in `namespaces` field.

All information cached in-memory. Any modification operation will rewrite file fully.

This kind of storage is limited by RAM size and requires quite a number of sycalls for update. So it's good for:

* Using data from another systems (JSON is a standard) or by humans
* Many reads and few writes

#### URL initialization

Do not forget to import package!

`file+json://<path>`

Where:

* `<path>` - path to file


## Usage

**Flat**

```go
stor := filestorage.NewFlat("path/to/directory")
// Close() not required but implemented as NOP
```


**Encoded**

```go
stor := filestorage.NewDefault("path/to/directory")
// Close() not required but implemented as NOP
```


{% include backend_tail.md page=page %}
