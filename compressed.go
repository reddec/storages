package storages

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

// Compressed storage where values are compressed by gzip
func Compressed(storage Storage) Storage {
	return &compressed{storage: storage}
}

type compressed struct {
	storage Storage
}

func (cs *compressed) Put(key []byte, data []byte) error {
	cdata, err := cs.packData(data)
	if err != nil {
		return err
	}
	return cs.storage.Put(key, cdata)
}

func (cs *compressed) Close() error {
	return cs.storage.Close()
}

func (cs *compressed) Get(key []byte) ([]byte, error) {
	cdata, err := cs.storage.Get(key)
	if err != nil {
		return nil, err
	}
	return cs.unpackData(cdata)
}

func (cs *compressed) Del(key []byte) error {
	return cs.storage.Del(key)
}

func (cs *compressed) Keys(handler func(key []byte) error) error {
	return cs.storage.Keys(handler)
}

func (cs *compressed) packData(data []byte) ([]byte, error) {
	buf := &bytes.Buffer{}
	writer := gzip.NewWriter(buf)
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (cs *compressed) unpackData(cdata []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(cdata))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(reader)
}
