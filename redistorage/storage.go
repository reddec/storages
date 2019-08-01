package redistorage

import (
	"github.com/go-redis/redis"
	"os"
	"storages"
)

type redisStorage struct {
	client *redis.Client
	key    string
}

func (rs *redisStorage) Put(key []byte, data []byte) error {
	return rs.client.HSet(rs.key, string(key), data).Err()
}

func (rs *redisStorage) Get(key []byte) ([]byte, error) {
	cmd := rs.client.HGet(rs.key, string(key))
	if cmd.Err() == redis.Nil {
		return nil, os.ErrNotExist
	}
	return cmd.Bytes()
}

func (rs *redisStorage) Del(key []byte) error {
	return rs.client.HDel(rs.key, string(key)).Err()
}

func (rs *redisStorage) Keys(handler func(key []byte) error) error {
	cmd := rs.client.HKeys(rs.key)
	if cmd.Err() == redis.Nil {
		return nil
	} else if cmd.Err() != nil {
		return cmd.Err()
	}
	keys, err := cmd.Result()
	if err != nil {
		return err
	}
	for _, k := range keys {
		err = handler([]byte(k))
		if err != nil {
			return err
		}
	}
	return nil
}

func (rs *redisStorage) Close() error { return rs.client.Close() }

// New storage wrapper around REDIS hashmap. Namespace is a hashkey
func NewClient(namespace string, client *redis.Client) storages.Storage {
	return &redisStorage{
		key:    namespace,
		client: client,
	}
}

// New REDIS client and storage wrapper
func New(namespace string, url string) (storages.Storage, error) {
	params, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return NewClient(namespace, redis.NewClient(params)), nil
}

// New REDIS client and storage wrapper. If URL is invalid - panic
func MustNew(namespace string, url string) storages.Storage {
	st, err := New(namespace, url)
	if err != nil {
		panic(err)
	}
	return st
}
