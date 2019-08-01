package queues

import (
	"reddec/storages"
	"sync"
)

type limitedQueue struct {
	queue storages.Queue
	lock  sync.Mutex
	limit int64
}

func (lq *limitedQueue) Peek() (id int64, data []byte, err error) {
	lq.lock.Lock()
	defer lq.lock.Unlock()
	return lq.queue.Peek()
}

func (lq *limitedQueue) Clean(end int64) error {
	lq.lock.Lock()
	defer lq.lock.Unlock()
	return lq.queue.Clean(end)
}

func (lq *limitedQueue) Size() int64 { return lq.queue.Size() }

func (lq *limitedQueue) First() int64 { return lq.queue.First() }

func (lq *limitedQueue) Last() int64 { return lq.queue.Last() }

func (lq *limitedQueue) Available() int64 {
	lq.lock.Lock()
	defer lq.lock.Unlock()
	return lq.limit - lq.Size()
}

func (lq *limitedQueue) Put(data []byte) (id int64, err error) {
	lq.lock.Lock()
	defer lq.lock.Unlock()
	id, err = lq.queue.Put(data)
	if err != nil {
		return
	}
	overflow := lq.Size() - lq.limit
	if overflow > 0 {
		// suppress errors during cleaning
		_ = lq.queue.Clean(lq.queue.First() + overflow)
	}
	return
}

func (lq *limitedQueue) Limit() int64 { return lq.limit }

func (lq *limitedQueue) Iterate(from int64) storages.Iterator { return lq.queue.Iterate(from) }

// Wrap any queue as limited queue
func Limited(queue storages.Queue, limit int64) storages.LimitedQueue {
	return &limitedQueue{
		limit: limit,
		queue: queue,
	}
}

// Create simple queue and wrap it to limited queue
func SimpleLimited(storage storages.Storage, limit int64) (storages.LimitedQueue, error) {
	q, err := Simple(storage)
	if err != nil {
		return nil, err
	}
	return Limited(q, limit), nil
}
