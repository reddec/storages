# Queues


    | data 1 | data 2 | data 3 | ... | data N |
    <-- tail (oldest)         head (latest) -->

Wrappers around KV-storage that makes a queues. Idea is to keep minimal and maximum id and use sequence to generate 
next key for KV storage.

Data always append to head and reads from tail (FIFO).

import: `github.com/reddec/storages/queues`

## Basic queue

* peek - get oldest data but do not remove
* put - push data to the end of queue
* discard - remove data oldest data

Constructors:

* `Naive`