---
backend: "In-Memory"
headline: "HashMap-based in-memory storage"
features: ["batch_writer", "namespace"]
---
# Memory DB

[![API docs](https://godoc.org/github.com/reddec/storages/memstorage?status.svg)](http://godoc.org/github.com/reddec/storages/memstorage)

* **import:**  `github.com/reddec/storages/memstorage`

Based on hashmap and RWLock in-memory storage. Values and keys are copied before put

For namespaces used Go `sync.Map`.

## Usage

```go
storage := New()
```

`Close()` is not required, however it is implemented as NOP.

## Features

{% for feature in features %}
{% include feature_{{feature}}.md %}
{% endfor %}
