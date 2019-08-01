package leveldbstorage

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"os"
	"storages"
)

type leveldbMap struct {
	db *leveldb.DB
}

func (bdp *leveldbMap) Put(key []byte, value []byte) error {
	return bdp.db.Put(key, value, nil)
}

func (bdp *leveldbMap) Get(key []byte) ([]byte, error) {
	data, err := bdp.db.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return nil, os.ErrNotExist
	}
	return data, err
}

func (bdp *leveldbMap) Del(key []byte) error {
	return bdp.db.Delete(key, nil)
}

func (bdp *leveldbMap) Keys(handler func(key []byte) error) error {
	it := bdp.db.NewIterator(nil, nil)
	defer it.Release()
	if it.Error() != nil {
		return it.Error()
	}
	for it.Next() {
		if it.Error() != nil {
			return it.Error()
		}
		err := handler(it.Key())
		if err != nil {
			return err
		}
	}
	return nil
}
func (bdp *leveldbMap) Close() error { return bdp.db.Close() }

// New storage, base on go-leveldb store
func New(location string) (storages.BatchedStorage, error) {
	db, err := leveldb.OpenFile(location, nil)
	if err != nil {
		return nil, err
	}
	return &leveldbMap{db: db}, nil
}

type dbBatch struct {
	batch *leveldb.Batch
	db    *leveldb.DB
}

func (bdp *leveldbMap) BatchWriter() storages.Writer {
	batch := new(leveldb.Batch)
	return &dbBatch{batch: batch, db: bdp.db}
}

func (dbt *dbBatch) Put(key []byte, data []byte) error {
	dbt.batch.Put(key, data)
	return nil
}

func (dbt *dbBatch) Close() error {
	return dbt.db.Write(dbt.batch, &opt.WriteOptions{})
}
