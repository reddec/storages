package rest

import (
	"bytes"
	"encoding/base64"
	"github.com/pkg/errors"
	"github.com/reddec/storages/std/memstorage"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewServer(t *testing.T) {
	back := memstorage.New()

	_ = back.Put([]byte("alice"), []byte("hello"))
	_ = back.Put([]byte("bob"), []byte("world"))

	server := httptest.NewServer(NewServer(back))
	defer server.Close()
	client := server.Client()
	// get keys
	res, err := client.Get(server.URL)
	if err != nil {
		t.Error("GET keys", err)
		return
	}
	if res.StatusCode != http.StatusOK {
		t.Error("bad status:", res.Status)
		return
	}
	keysData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error("GET keys/decode", err)
		return
	}
	keysChunks := strings.Split(string(keysData), "\n")
	if len(keysChunks) != 2 {
		t.Error("num keys miss-match")
		return
	}
	var keyAlice bool
	var keyBob bool
	for _, line := range keysChunks {
		key, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			t.Error("decode key:", err)
			return
		}
		if string(key) == "alice" {
			keyAlice = true
		} else if string(key) == "bob" {
			keyBob = true
		} else {
			t.Error("unknown key:", string(key))
			return
		}
	}
	if !keyAlice || !keyBob {
		t.Error("no all keys arrived")
		return
	}
	// get key - alice
	dataAlice, err := fetchKey("alice", server)
	if err != nil {
		t.Error("get alice", err)
		return
	}
	if string(dataAlice) != "hello" {
		t.Error("alice has a wrong value")
		return
	}
	// get key - bob
	dataBob, err := fetchKey("bob", server)
	if err != nil {
		t.Error("get bob", err)
		return
	}
	if string(dataBob) != "world" {
		t.Error("bob has a wrong value")
		return
	}
	// put carl
	res, err = client.Post(server.URL+"/"+base64.StdEncoding.EncodeToString([]byte("carl")), "", bytes.NewBufferString("hell in world"))
	if err != nil {
		t.Error("put carl", err)
		return
	}
	if res.StatusCode != http.StatusNoContent {
		t.Error("bad status:", res.Status)
		return
	}
	// check carl
	dataCarl, err := fetchKey("carl", server)
	if err != nil {
		t.Error("get carl", err)
		return
	}
	if string(dataCarl) != "hell in world" {
		t.Error("carl has a wrong value:", string(dataCarl))
		return
	}
	// delete alice
	rq, err := http.NewRequest(http.MethodDelete, server.URL+"/"+base64.StdEncoding.EncodeToString([]byte("alice")), nil)
	if err != nil {
		t.Error("request:", err)
		return
	}
	res, err = client.Do(rq)
	if err != nil {
		t.Error("remove alice", err)
		return
	}
	if res.StatusCode != http.StatusNoContent {
		t.Error("bad status:", res.Status)
		return
	}
	// check alice
	dataAlice, err = fetchKey("alice", server)
	if err == nil {
		t.Error("alice should not exists but got", string(dataAlice))
		return
	}
	if !strings.Contains(err.Error(), "Not Found") {
		t.Error(err.Error(), "but should be not found")
		return
	}
}

func fetchKey(key string, server *httptest.Server) ([]byte, error) {
	client := server.Client()
	res, err := client.Get(server.URL + "/" + base64.StdEncoding.EncodeToString([]byte(key)))
	if err != nil {
		return nil, errors.Wrapf(err, "get key: %s", key)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("get key: %s, status: %s", key, res.Status)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "get key: %s, read data", key)
	}
	return data, nil
}
