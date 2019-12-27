package graph

// Node key
type NodeKey []byte

type Node interface {
	// node payload
	Data() ([]byte, error)
	// node id (key)
	Key() NodeKey
	// list of linked keys to the node
	Linked() ([]NodeKey, error)
	// update node payload
	SetData(value []byte) error
	// update node links
	SetLinked(keys []NodeKey) error
	// open another node by key using current node implementation of graph
	Open(key NodeKey) Node
}
