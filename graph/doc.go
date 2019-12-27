// Graph-like operations on top of storage
/*
The package is in work-in-progress status.

Classical data structure where each node (element) addressed by unique key (id) and should have only one value
and multiple links: keys of another nodes. Such data structure should well reflect real world relationship like so:
social network (person can have friends),  users accounts (several accounts for one user) and so on.

                NODE                    NODE
             +--------+              +--------+
    node id  |  key   |     +------> |  key   |
             +--------+     |        | ....   |
    payload  | value  |     |        +--------+
             +--------+     |
             |  key0  |  ---+
    links    |  key1  |
	         |  ...   |
             +--------+


Graph may contain only forward link as well as backwards links also. It depends of implementation.
*/
package graph
