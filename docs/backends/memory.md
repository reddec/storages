---
backend: "In-Memory"
package: "std/memstorage"
headline: "HashMap-based in-memory storage"
features: ["batch_writer", "namespace", "clearable"]
project_url: ""
---
{% include backend_head.md page=page %}

Based on hashmap and RWLock in-memory storage. Values and keys are copied before put

For namespaces used Go `sync.Map`.

### URL initialization

Do not forget to import package!

`memory://`

## Usage

```go
storage := memstorage.New()
```

`Close()` is not required, however it is implemented as NOP.

{% include backend_tail.md page=page %}
