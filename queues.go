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
	// iterate over items from first (if from is 0) to last (next should be called first) or till first error.
	// Iterator keeps min and max sequence number so cleaning items during iteration may cause iteration stop
	Iterate(from int64) Iterator
}

// Queue iterator
type Iterator interface {
	// Is queue has next value
	Next() bool
	// Current id
	ID() int64
	// Current value
	Value() []byte
}

// Queue with limited size
type LimitedQueue interface {
	Queue
	// available space (limit - size)
	Available() int64
	// limit
	Limit() int64
}
