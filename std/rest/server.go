package rest

import (
	"encoding/base64"
	"github.com/reddec/storages"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Creates new http handler and provides REST-like access to storage.
//
// GET / - array of all keys. Each key - base64 encoded. New line - new key. Stream is chunk encoded. Returns 200
//
// GET /:key - content of key. Returns 404 if key not found. key should be base64 encoded
//
// POST,PUT,PATCH /:key - update or insert value for key. Returns 204 on success. key should be base64 encoded
//
// DELETE /:key - remove key. Returns 204 on success. key should be base64 encoded
func NewServer(backed storages.Storage) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.URL.Path == "/" {
			if r.Method == http.MethodGet {
				listKeys(backed, w, r)
				return
			} else {
				http.Error(w, "no method", http.StatusMethodNotAllowed)
				return
			}
		}
		key, err := base64.StdEncoding.DecodeString(r.URL.Path[1:])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			getKey(key, backed, w, r)
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			postKey(key, backed, w, r)
		case http.MethodDelete:
			removeKey(key, backed, w, r)
		default:
			http.Error(w, "unknown operation", http.StatusMethodNotAllowed)
		}
	})

	return mux
}

func listKeys(backed storages.Storage, w http.ResponseWriter, r *http.Request) {
	var sent bool
	err := backed.Keys(func(key []byte) error {
		text := base64.StdEncoding.EncodeToString(key)
		if !sent {
			w.Header().Set("Content-Encoding", "base64")
			w.WriteHeader(http.StatusOK)
		} else {
			text = "\n" + text
		}
		sent = true

		_, err := w.Write([]byte(text))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if sent {
			log.Println("[ERROR]", err)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func getKey(key []byte, backed storages.Storage, w http.ResponseWriter, r *http.Request) {
	data, err := backed.Get(key)
	if err == os.ErrNotExist {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func postKey(key []byte, backed storages.Storage, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = backed.Put(key, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func removeKey(key []byte, backed storages.Storage, w http.ResponseWriter, r *http.Request) {
	err := backed.Del(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
