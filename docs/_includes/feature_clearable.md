### Clearable

Support [Clearable](https://godoc.org/github.com/reddec/storages#Clearable) interface.

Allows clean all internal information by one operation.

**Example:**
  
```go
err := storage.Clean()
if err != nil {
    panic("failed to clean")
}
```