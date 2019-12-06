# Redundancy


Redundancy is a process to copy data through multiple storages. 

                  storage1 (data from 0..N)
                🡕
    data(0...N) 🡒 storage1 (data from 0..N)
                🡖 
                  storage1 (data from 0..N)
                  
In contrast of [sharding](sharding) where data scattered over multiple storages without duplicates, the
redundancy oppositely coping data to all underlying storages (depends of writer strategy).

Briefly, sharding is for horizontal scaling, redundancy is for reliability. 

Redundant storage mimics to usual [Storage](https://godoc.org/github.com/reddec/storages#Storage) interface so 
target consumers should work with it as usual.

General constructor for redundant storage is [Redundant(writerStrategy, readerStrategy, deduplication, ...storages)](https://godoc.org/github.com/reddec/storages#Redundant).

Shorthand constructor [RedundantAll(deduplication, ...storages)](https://godoc.org/github.com/reddec/storages#RedundantAll) 
offers those default strategy:

* All storages should successfully be written (strategy [AtLeast](https://godoc.org/github.com/reddec/storages#AtLeast))
* First non-empty, non-error result will be return (strategy [First](https://godoc.org/github.com/reddec/storages#First))