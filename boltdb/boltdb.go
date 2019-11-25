package boltdb

import (
	"github.com/reddec/storages"
	"go.etcd.io/bbolt"
	"os"
)

const (
	DefaultBucket = "DEFAULT" // default name for bucket
)

func New(location string, namespace []byte) (*boltDB, error) {
	db, err := bbolt.Open(location, 0755, nil)
	if err != nil {
		return nil, err
	}
	return &boltDB{
		db:     db,
		bucket: namespace,
	}, nil
}

func NewDefault(location string) (*boltDB, error) {
	return New(location, []byte(DefaultBucket))
}

type boltDB struct {
	db     *bbolt.DB
	bucket []byte
	nested bool
}

func (bdb *boltDB) DelNamespace(name []byte) error {
	return bdb.db.Update(func(tx *bbolt.Tx) error {
		return tx.DeleteBucket(name)
	})
}

func (bdb *boltDB) Put(key []byte, data []byte) error {
	return bdb.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bdb.bucket)
		if err != nil {
			return err
		}
		return bucket.Put(key, data)
	})
}

func (bdb *boltDB) Close() error {
	if bdb.nested {
		return nil
	}
	return bdb.db.Close()
}

func (bdb *boltDB) Get(key []byte) ([]byte, error) {
	var ans []byte
	err := bdb.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bdb.bucket)
		if bucket == nil {
			return os.ErrNotExist
		}
		value := bucket.Get(key)
		if value == nil {
			return os.ErrNotExist
		}
		ans = value
		return nil
	})
	return ans, err
}

func (bdb *boltDB) Del(key []byte) error {
	return bdb.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bdb.bucket)
		if bucket == nil {
			return nil
		}
		return bucket.Delete(key)
	})
}

func (bdb *boltDB) Keys(handler func(key []byte) error) error {
	return bdb.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bdb.bucket)
		if bucket == nil {
			return nil
		}
		return bucket.ForEach(func(k, v []byte) error {
			return handler(k)
		})
	})
}

func (bdb *boltDB) Namespace(name []byte) (storages.Storage, error) {
	err := bdb.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(name)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &boltDB{
		db:     bdb.db,
		bucket: name,
		nested: true,
	}, nil
}

func (bdb *boltDB) Namespaces(handler func(name []byte) error) error {
	return bdb.db.View(func(tx *bbolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
			return handler(name)
		})
	})
}
