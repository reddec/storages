# Queues

Wrappers around KV-storage that makes a queues. Idea is to keep minimal and maximum id and use sequence to generate 
next key for KV storage.

import: `github.com/reddec/storages/queues`

## Basic queue

* peek - get last but do not remove
* put - push data to the end of queue
* clean - remove data from the first till specified sequence id. Remove all is: `Clean(queue.Last()+1)`

Constructors:

* `Simple`
* `SimpleBounded`

## Limited queue

Queue that removes old items if no more space available (like circular buffer) on `Put` operation.

Constructors:

* `Limited`
* `SimpleLimited` - shorthand for `Limited(Simple(...))`
