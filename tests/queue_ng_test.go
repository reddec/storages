package tests

import (
	"github.com/reddec/storages"
	"github.com/reddec/storages/queues"
	"github.com/reddec/storages/std/memstorage"
	"os"
	"testing"
)

func TestQueueNQ(t *testing.T) {
	mem := memstorage.New()
	defer mem.Close()
	testQueue(func() (queue storages.Queue, err error) {
		return queues.NaiveQueue(mem)
	}, t)
}

func testQueue(queueFactory func() (storages.Queue, error), t *testing.T) {
	q, err := queueFactory()
	if err != nil {
		t.Error("open", err)
		return
	}
	// empty queue
	_, err = q.Get()
	if err != os.ErrNotExist {
		t.Error("empty queue contains something?")
		return
	}
	err = q.Discard()
	if err != os.ErrNotExist {
		t.Error("empty queue contains something during discard?")
		return
	}
	_, err = q.Peek()
	if err != os.ErrNotExist {
		t.Error("empty queue contains something during peek?")
		return
	}
	// add to queue
	err = q.Put([]byte("alice"))
	if err != nil {
		t.Error("put", err)
		return
	}
	err = q.Put([]byte("bob"))
	if err != nil {
		t.Error("put", err)
		return
	}
	err = q.Put([]byte("clark"))
	if err != nil {
		t.Error("put", err)
		return
	}
	// peek
	data, err := q.Peek()
	if err != nil {
		t.Error("peek", err)
		return
	}
	if string(data) != "alice" {
		t.Error("where is alice")
		return
	}
	// repeat peek - should be the same
	data, err = q.Peek()
	if err != nil {
		t.Error("peek", err)
		return
	}
	if string(data) != "alice" {
		t.Error("where is alice")
		return
	}
	// now get one
	data, err = q.Get()
	if err != nil {
		t.Error("get", err)
		return
	}
	if string(data) != "alice" {
		t.Error("where is alice")
		return
	}
	// next one should be different
	data, err = q.Get()
	if err != nil {
		t.Error("get", err)
		return
	}
	if string(data) != "bob" {
		t.Error("where is bob")
		return
	}
	// peek and then discard
	data, err = q.Peek()
	if err != nil {
		t.Error("peek", err)
		return
	}
	if string(data) != "clark" {
		t.Error("where is clark")
		return
	}
	err = q.Discard()
	if err != nil {
		t.Error("discard", err)
		return
	}

	// should be empty
	_, err = q.Peek()
	if err != os.ErrNotExist {
		t.Error("empty queue contains something during peek?")
		return
	}

	// add something and recreate queue
	err = q.Put([]byte("alice"))
	if err != nil {
		t.Error("put", err)
		return
	}

	q, err = queueFactory()
	if err != nil {
		t.Error("open", err)
		return
	}
	// check saved
	data, err = q.Peek()
	if err != nil {
		t.Error("peek", err)
		return
	}
	if string(data) != "alice" {
		t.Error("where is alice")
		return
	}
}
