package config

import (
	"encoding/json"
	"github.com/reddec/storages"
	"github.com/reddec/storages/std/memstorage"
	"reflect"
	"strings"
	"testing"
)

func TestParseJSON(t *testing.T) {
	var stor = memstorage.New()
	setConfig(stor, t, "main", Sharded{Shards: []string{"shard1", "shard2", "shard3"}})
	setConfig(stor, t, "shard1", Redundant{Storages: []string{"data1", "data2"}})
	setConfig(stor, t, "shard2", Redundant{Storages: []string{"data2", "data3"}})
	setConfig(stor, t, "shard3", Redundant{Storages: []string{"data3", "data1"}})
	setConfig(stor, t, "data1", Simple{URL: "memory://"})
	setConfig(stor, t, "data2", Simple{URL: "memory://"})
	setConfig(stor, t, "data3", Simple{URL: "memory://"})
	setConfig(stor, t, "data4", Simple{URL: "memory://"})

	root, err := ParseJSON([]byte("main"), stor)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()

}

func setConfig(stor storages.Storage, t *testing.T, key string, config interface{}) {

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	var iconfig map[string]interface{}
	err = json.Unmarshal(data, &iconfig)
	if err != nil {
		t.Fatal(err)
	}
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	iconfig["kind"] = strings.ToLower(v.Type().Name())

	data, err = json.Marshal(iconfig)
	if err != nil {
		t.Fatal(err)
	}

	err = stor.Put([]byte(key), data)
	if err != nil {
		t.Fatal(err)
	}
}
