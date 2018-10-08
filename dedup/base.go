package dedup

// Deduplicate primitive: check if key is already saved and save key
type Dedup interface {
	// Is key already save?
	IsDuplicated(key []byte) (bool, error)
	// Save key for future checks
	Save(key []byte) error
}
