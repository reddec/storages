---
backend: "Mock"
package: "std/memstorage"
headline: "Mocking storage that do nothing"
features: ["batch_writer"]
project_url: ""
---
{% include backend_head.md page=page %}

No-Operation storage that drops any content and returns not-exists on any request.

Useful for mocking, performance testing or for any other logic that needs discard storage.

## Usage

**Example**

```go
stor := NewNOP()
// no need to call Close()
```


{% include backend_tail.md page=page %}

