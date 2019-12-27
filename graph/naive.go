package graph

import (
	"bytes"
	"encoding/gob"
	"github.com/reddec/storages"
)

// Naive graph stores value and links in one storage cell with GOB encoding (could be changed in a future)
func Naive(storage storages.Storage, key NodeKey) *graphNode {
	return &graphNode{
		key:     key,
		storage: storage,
	}
}

type cell struct {
	Links []NodeKey
	Value []byte
}

type graphNode struct {
	cell    *cell
	key     NodeKey
	storage storages.Storage
}

func (gn *graphNode) Data() ([]byte, error) {
	if gn.cell != nil {
		return gn.cell.Value, nil
	}
	v, err := gn.storage.Get(gn.key)
	if err != nil {
		return nil, err
	}

	err = gob.NewDecoder(bytes.NewBuffer(v)).Decode(gn.cell)
	if err != nil {
		gn.cell = nil
		return nil, err
	}
	return gn.cell.Value, nil
}

func (gn *graphNode) Key() NodeKey {
	return gn.key
}

func (gn *graphNode) Linked() ([]NodeKey, error) {
	_, err := gn.Data()
	if err != nil {
		return nil, err
	}
	return gn.cell.Links, nil
}

func (gn *graphNode) SetData(value []byte) error {
	var keys []NodeKey
	if gn.cell != nil {
		keys = gn.cell.Links
	}
	return gn.update(value, keys)
}

func (gn *graphNode) SetLinked(keys []NodeKey) error {
	var data []byte
	if gn.cell != nil {
		data = gn.cell.Value
	}
	return gn.update(data, keys)
}

func (gn *graphNode) Open(key NodeKey) Node { return Naive(gn.storage, key) }

func (gn *graphNode) update(value []byte, linked []NodeKey) error {
	buf := &bytes.Buffer{}

	cell := &cell{
		Value: value,
		Links: linked,
	}

	err := gob.NewEncoder(buf).Encode(cell)
	if err != nil {
		return err
	}
	gn.cell = cell
	return nil
}
