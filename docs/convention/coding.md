# Coding style

* V0 - return abstract interfaces
* V1 - follow 'accept interfaces, return structs'

Since V1 all implementations should return non-exported reference to structure (see `boltdb` wrapper as an example). Standard wrappers will be replace as sooner as possible, 
however it should not affect code that already using current library.