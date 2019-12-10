package tests

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/pkg/errors"
	"github.com/reddec/storages"
	"github.com/reddec/storages/dedup"
	"github.com/reddec/storages/sharded"
	"github.com/reddec/storages/std/awsstorage"
	"github.com/reddec/storages/std/boltdb"
	"github.com/reddec/storages/std/filestorage"
	"github.com/reddec/storages/std/leveldbstorage"
	"github.com/reddec/storages/std/memstorage"
	"github.com/reddec/storages/std/redistorage"
	"github.com/reddec/storages/std/rest"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func Test_Storages(t *testing.T) {
	var testDir string
	var stor storages.Storage
	var err error

	err = os.RemoveAll("../test")
	if err != nil {
		t.Error(err)
		return
	}

	err = os.MkdirAll("../test", 0755)
	if err != nil {
		t.Error(err)
		return
	}

	// file storage
	testDir = "../test/file-storage"
	stor = filestorage.NewDefault(testDir)
	testStorage(t, stor, testDir, true)

	// level db storage
	testDir = "../test/leveldb-storage"
	stor, err = leveldbstorage.New(testDir)
	if err != nil {
		t.Fatal("leveldb:", err)
	}
	testStorage(t, stor, testDir, true)

	TestMemory(t)
	// redis storage (REDIS should be installed and started on default port)
	testDir = "../test/redis-storage"
	stor = redistorage.MustNew("data", "redis://127.0.0.1")
	testShouldBeNS(t, stor)
	testStorage(t, stor, testDir, true)

	// AWS storage
	//TestAWS(t)
	// Flat files
	TestFlat(t)
	TestFileJSON(t)
	// Test bolt
	TestBolt(t)
	// Test REST client and server
	TestRest(t)
}

func TestRest(t *testing.T) {
	back := memstorage.New()
	server := httptest.NewServer(rest.NewServer(back))
	defer server.Close()

	stor := rest.NewClient(server.URL)
	defer stor.Close()
	testStorage(t, stor, "", true)
}
func TestRedundancy(t *testing.T) {
	var used []storages.Storage

	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			used = append(used, memstorage.NewNOP()) // simulate successful write but broken read
		} else {
			used = append(used, memstorage.New())
		}
	}

	rdr := storages.RedundantAll(dedup.Offloaded(memstorage.New()), used...)
	defer rdr.Close()

	testStorage(t, rdr, "", true)
}

func TestShard(t *testing.T) {
	var used []storages.Storage
	pool := sharded.NewHashed(3, func(shardID uint32) (storage storages.Storage, e error) {
		mem := memstorage.New()
		used = append(used, mem)
		return mem, nil
	})

	shard := storages.Sharded(pool)

	testStorage(t, shard, "", true)

	if err := shard.Put([]byte("alice"), []byte("1")); err != nil {
		t.Error(err)
		return
	}
	if err := shard.Put([]byte("bob1"), []byte("1")); err != nil {
		t.Error(err)
		return
	}
	if err := shard.Put([]byte("carly"), []byte("1")); err != nil {
		t.Error(err)
		return
	}
	first, _ := storages.AllKeysString(used[0])
	second, _ := storages.AllKeysString(used[1])
	third, _ := storages.AllKeysString(used[2])
	assert.Len(t, first, 1)
	assert.Len(t, second, 1)
	assert.Len(t, third, 1)
	t.Logf("1st: %+v, 2nd: %+v, 3d: %+v", first, second, third)
}

func TestBolt(t *testing.T) {
	testFile := "../test/boltd.db"
	stor, err := boltdb.NewDefault(testFile)
	if err != nil {
		t.Error(err)
		return
	}
	defer stor.Close()
	testShouldBeNS(t, stor)
	testStorage(t, stor, "", true)
}

func TestMemory(t *testing.T) {
	// memory storage
	testDir := "../test/memory-storage"
	stor := memstorage.New()
	testShouldBeNS(t, stor)
	testStorage(t, stor, testDir, true)
}

func TestFlat(t *testing.T) {
	testDir := "../test/flat-file-storage"
	stor := filestorage.NewFlat(testDir)
	testShouldBeNS(t, stor)
	testStorage(t, stor, testDir, true)
}

func TestFileJSON(t *testing.T) {
	testFile := "../test/data.json"
	stor := filestorage.NewJSONFile(testFile)
	testShouldBeNS(t, stor)
	testStorage(t, stor, "", true)
}

func TestAWS(t *testing.T) {
	config := aws.NewConfig()
	config.Credentials = credentials.NewEnvCredentials()
	stor, err := awsstorage.New(os.Getenv("BUCKET"), config)
	if err != nil {
		t.Error(err)
		return
	}
	defer stor.Close()
	testStorage(t, stor, "", true)
}

func testShouldBeNS(t *testing.T, storage storages.Storage) {
	if _, ok := storage.(storages.NamespacedStorage); !ok {
		t.Errorf("%v should be namespaced storage", reflect.ValueOf(storage).Elem().Type().Name())
	}
}

func testStorage(t *testing.T, storage storages.Storage, testDir string, testNested bool) {
	if testDir != "" {
		os.RemoveAll(testDir)
		err := os.MkdirAll(testDir, 0755)
		if err != nil {
			t.Fatal("init test:", err)
		}
	}
	t.Log("Testing", reflect.ValueOf(storage).Elem().Type().Name())

	// pre check
	err := storage.Keys(func(key []byte) error {
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
	if err != nil {
		t.Error("get test1:", err)
		return
	}
	if string(data) != "hello world 1" {
		t.Error("corrupted value for test1:", string(data))
		return
	}
	data, err = storage.Get([]byte("test2"))
	if err != nil {
		t.Error("get test2:", err)
		return
	}
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

	if ns, ok := storage.(storages.NamespacedStorage); ok && testNested {
		t.Log("testing namespaces")
		testNamespaces(ns, t, testNested)
	}
}

func testNamespaces(storage storages.NamespacedStorage, t *testing.T, testNested bool) {
	err := storage.Namespaces(func(name []byte) error {
		t.Log("warning! Already exists namespace in empty storage:", string(name))
		return nil
	})
	if err != nil {
		t.Error(err)
		return
	}

	ns, err := storage.Namespace([]byte("test1"))
	if err != nil {
		t.Error(err)
		return
	}
	if testNested {
		testStorage(t, ns, "", false)
	}
	err = ns.Put([]byte("A"), []byte("B"))
	if err != nil {
		t.Error(err)
		return
	}
	list, err := storages.AllNamespacesString(storage)
	if err != nil {
		t.Error(err)
		return
	}
	if len(list) == 0 {
		t.Error("no ns created")
		return
	}
	var found bool
	for _, k := range list {
		if k == "test1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("created namespace not found")
		return
	}
	err = storage.DelNamespace([]byte("test1"))
	if err != nil {
		t.Error(err)
		return
	}
	// check removed namespace
	list, err = storages.AllNamespacesString(storage)
	if err != nil {
		t.Error(err)
		return
	}
	found = false
	for _, k := range list {
		if k == "test1" {
			found = true
			break
		}
	}
	if found {
		t.Error("removed namespace still exists")
		return
	}
}
