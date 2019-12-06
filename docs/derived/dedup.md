
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


## Naive

Naive implementation of deduplicate process: simply keep keys as-is, remove old keys when amount (quantity) increased up to
`maxKeys * cleanFactor` till `maxKeys count`.

Relay on storages to detect order of keys.

Cleaning of old keys initiates in `Save()`` method automatically in a same thread.

Properties:

* `maxKeys` - maximum keys to store after cleanup
* `cleanFactor` - multiply factor of `maxKeys` that triggers cleanup process


## Offloaded

Offloaded deduplication is a wrapper around storage that checks and store keys with random unique iteration id.
In case if storage does not support `Clearable` interface, unique random iteration id could be used for distinguish
keys from different iterations. **Good for large data set in case wrapped storage is a disk-based solution.**

Offloaded deduplication supports `Clearable` interface by it self. Within `Clean` operation Offloaded deduplication
 instance tries to Clear underlying storage and resets iteration id.  

{% include feature_clearable.md %}