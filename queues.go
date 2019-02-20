package storages

// Wrapper around storage that makes sequential data inserting and peeking
// Without external access to the data Queue guarantees that sequences are without space and strictly increasing.
// If queue has no data sequence is flushed (starts from 0) after restart
type Queue interface {
	// put data to queue using new sequence id
	Put(data []byte) (id int64, err error)
	// peek latest data
	Peek() (id int64, data []byte, err error)
	// clean data in range: [first;end)
	Clean(end int64) error
	// size of queue (last-first)
	Size() int64
	// first (oldest) sequence id
	First() int64
	// last (latest) sequence id
	Last() int64
}
