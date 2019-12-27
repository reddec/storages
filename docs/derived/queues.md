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


## HTTP expose

It's possible to expose queue over HTTP by `NewServer(queue)`  


| Method   | Path   | Success status | Description |
|----------|--------|----------------|-------------
| `GET`    | `/`    | 200            | Peek last message in queue (404 NotFound if queue is empty)
| `POST`   | `/`    | 204            | Add message to queue
| `DELETE` | `/`    | 200            | Get last message from queue and remove it. Last message will be returned otherwise 404 not found
