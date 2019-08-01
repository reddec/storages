package tests

import (
	"fmt"
	"storages"
	"storages/memstorage"
	"storages/queues"
	"strconv"
	"testing"
)

func TestQueue(t *testing.T) {
	mem := memstorage.New()
	queue, err := queues.Simple(mem)
	if err != nil {
		t.Error(err)
		return
	}
	id, err := queue.Put([]byte("hello"))
	if err != nil {
		t.Error(err)
		return
	}
	if id != 0 {
		t.Error("invalid id:", id)
		return
	}
	id2, err := queue.Put([]byte("world"))
	if err != nil {
		t.Error(err)
		return
	}
	if id2 != 1 {
		t.Error("invalid id:", id)
		return
	}
	// now check data
	id, data, err := queue.Peek()
	if err != nil {
		t.Error(err)
		return
	}
	if id != 1 {
		t.Error("invalid peek id:", id)
		return
	}
	if string(data) != "world" {
		t.Error("different data on top of queue:", string(data))
		return
	}

	if f := queue.First(); f != 0 {
		t.Error("different first value:", f)
		return
	}
	if f := queue.Last(); f != 1 {
		t.Error("different last value:", f)
		return
	}
	if f := queue.Size(); f != 2 {
		t.Error("different size:", f)
		return
	}
	if !validateHelloWorld(t, queue) {
		return
	}
	// test recovery
	q2, err := queues.Simple(mem)
	if err != nil {
		t.Error(err)
		return
	}
	if !validateHelloWorld(t, q2) {
		return
	}
	it := queue.Iterate(0)
	for it.Next() {
		fmt.Println("ID:", it.ID(), "Value:", string(it.Value()))
	}
	// test remove one
	err = queue.Clean(queue.Last()) // up to
	if err != nil {
		t.Error(err)
		return
	}
	if queue.Size() != 1 {
		t.Error("queue should have size 1")
		return
	}
	err = queue.Clean(queue.Last() + 1)
	if err != nil {
		t.Error(err)
		return
	}
	if queue.Size() != 0 {
		t.Error("queue should have size 0")
		return
	}

}

func TestLimitedQueue(t *testing.T) {
	mem := memstorage.New()
	q, err := queues.SimpleLimited(mem, 2)
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < 10; i++ {
		_, err = q.Put([]byte(strconv.Itoa(i)))
		if err != nil {
			t.Error(err)
			return
		}
	}
	if q.Size() != 2 {
		t.Error("size should be 2")
		return
	}
	_, data, err := q.Peek()
	if err != nil {
		t.Error(err)
		return
	}
	if string(data) != "9" {
		t.Error("should be 9 but got: " + string(data))
		return
	}
}

func validateHelloWorld(t *testing.T, queue storages.Queue) bool {
	id, data, err := queue.Peek()
	if err != nil {
		t.Error(err)
		return false
	}
	if id != 1 {
		t.Error("invalid peek id:", id)
		return false
	}
	if string(data) != "world" {
		t.Error("different data on top of queue:", string(data))
		return false
	}

	if f := queue.First(); f != 0 {
		t.Error("different first value:", f)
		return false
	}
	if f := queue.Last(); f != 1 {
		t.Error("different last value:", f)
		return false
	}
	if f := queue.Size(); f != 2 {
		t.Error("different size:", f)
		return false
	}
	return true
}
