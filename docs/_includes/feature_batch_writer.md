### Batch writing

Support [BatchedStorage](https://godoc.org/github.com/reddec/storages#BatchedStorage) interface.

It allows to cache `Put` operations and execute them in one batch. 
In general case it may increase write throughput.

Batch implements [Writer](https://godoc.org/github.com/reddec/storages#Writer) interface.

**Example:**
  
```go
batchWriter := storage.Batch()

batchWriter.Put([]byte("key1"),[]byte("value1")
batchWriter.Put([]byte("key2"),[]byte("value2")
//...
batchWriter.Put([]byte("keyN"),[]byte("valueN")

// flush/write batch to storage
batchWriter.Close() 
```