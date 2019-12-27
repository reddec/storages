package queues

import (
	"github.com/reddec/storages"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

// Creates new http handler and provides REST-like access to queue.
//
// GET / - peek last message in queue (404 NotFound if queue is empty)
//
// POST,PUT / - add message to queue
//
// DELETE / - get last message from queue and remove it. Last message will be returned otherwise 404 not found
func NewServer(q storages.Queue) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()

		switch request.Method {
		case http.MethodGet: // peek last
			data, err := q.Peek()
			reply(data, err, request, writer)
		case http.MethodPost, http.MethodPut: // push to queue
			data, err := ioutil.ReadAll(request.Body)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
			err = q.Put(data)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			writer.WriteHeader(http.StatusNoContent)
		case http.MethodDelete: // get last and discard
			data, err := q.Get()
			reply(data, err, request, writer)
		}
	})
	return mux
}

func reply(data []byte, err error, request *http.Request, writer http.ResponseWriter) {
	if err == os.ErrNotExist {
		http.NotFound(writer, request)
		return
	} else if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/octet-stream")
	writer.Header().Set("Content-Length", strconv.Itoa(len(data)))
	writer.WriteHeader(http.StatusOK)
	writer.Write(data)
}
