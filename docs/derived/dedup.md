
# Collection of deduplicate methods

All implementations should follow those interface

```golang
// Deduplicate primitive: check if key is already saved and save key
type Dedup interface {
	// Is key already save?
	IsDuplicated(key []byte) (bool, error)
	// Save key for future checks
	Save(key []byte) (error)
}
```


### Naive


Properties:

* `maxKeys` - maximum keys to store after cleanup
* `cleanFactor` - multiply factor of `maxKeys` that triggers cleanup process

Naive implementation of deduplicate process: simply keep keys as-is, remove old keys when amount (quantity) increased up to
`maxKeys * cleanFactor` till `maxKeys count`.

Relay on storages to detect order of keys.

Cleaning of old keys initiates in `Save()`` method automatically in a same thread.
