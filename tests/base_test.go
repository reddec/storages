package tests

import (
	"github.com/pkg/errors"
	"github.com/reddec/storages"
	"github.com/reddec/storages/filestorage"
	"github.com/reddec/storages/leveldbstorage"
	"github.com/reddec/storages/memstorage"
	"github.com/reddec/storages/redistorage"
	"os"
	"testing"
)

func Test_Storages(t *testing.T) {
	var testDir string
	var stor storages.Storage
	var err error

	// file storage
	testDir = "../test/file-storage"
	stor = filestorage.NewDefault(testDir)
	testStorage(t, stor, testDir)

	// level db storage
	testDir = "../test/leveldb-storage"
	stor, err = leveldbstorage.New(testDir)
	if err != nil {
		t.Fatal("leveldb:", err)
	}
	testStorage(t, stor, testDir)

	// memory storage
	testDir = "../test/memory-storage"
	stor = memstorage.New()
	testStorage(t, stor, testDir)

	// redis storage (REDIS should be installed and started on default port)
	testDir = "../test/redis-storage"
	stor = redistorage.MustNew("data", "redis://127.0.0.1")
	testStorage(t, stor, testDir)
}

func testStorage(t *testing.T, storage storages.Storage, testDir string) {
	os.RemoveAll(testDir)
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatal("init test:", err)
	}

	// pre check
	err = storage.Keys(func(key []byte) error {
		return errors.New("should be no keys on in empty dir")
	})

	if err != nil {
		t.Error("get keys:", err)
		return
	}

	// add data
	err = storage.Put([]byte("test1"), []byte("hello world 1"))
	if err != nil {
		t.Error("put test1:", err)
		return
	}
	err = storage.Put([]byte("test2"), []byte("hello world 2"))
	if err != nil {
		t.Error("put test2:", err)
		return
	}
	// get data
	data, err := storage.Get([]byte("test1"))
	if string(data) != "hello world 1" {
		t.Error("corrupted value for test1:", string(data))
		return
	}
	data, err = storage.Get([]byte("test2"))
	if string(data) != "hello world 2" {
		t.Error("corrupted value for test2:", string(data))
		return
	}
	// check keys
	var test1, test2 bool
	err = storage.Keys(func(key []byte) error {
		s := string(key)
		switch s {
		case "test1":
			test1 = true
		case "test2":
			test2 = true
		default:
			return errors.New("unknown key found: " + s)
		}
		return nil
	})

	if err != nil {
		t.Error("failed iterate storage keys:", err)
		return
	}

	if !test1 {
		t.Error("test1 key not found")
		return
	}
	if !test2 {
		t.Error("test2 key not found")
	}

	// remove
	err = storage.Del([]byte("test1"))
	if err != nil {
		t.Error("del test1:", err)
		return
	}
	_, err = storage.Get([]byte("test1"))
	if err != os.ErrNotExist {
		t.Error("get removed key test1 caused NOT ErrNotExist error:", err)
		return
	}
	err = storage.Del([]byte("test2"))
	if err != nil {
		t.Error("del test2:", err)
		return
	}
	_, err = storage.Get([]byte("test2"))
	if err != os.ErrNotExist {
		t.Error("get removed key test2 caused NOT ErrNotExist error:", err)
		return
	}
}
