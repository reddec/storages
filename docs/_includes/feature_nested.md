* Supports `NestedStorage` interfaces.

  It allows make nested sub-storages with independent
  key space.
  
  Example:
  
```go
aliceStorage :=  storage.Namespace([]byte("alice"))
bobStorage :=  storage.Namespace([]byte("bob"))
```