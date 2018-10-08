package dedup

import (
	"github.com/reddec/storages"
	"github.com/reddec/storages/memstorage"
	"testing"
)

func TestNewNaive(t *testing.T) {
	mem := memstorage.New()
	defer mem.Close()
	nv, err := NewNaive(mem, 2, 2)
	if err != nil {
		t.Error("failed initialize naive dedup:", err)
		return
	}

	keys := [][]byte{[]byte("01"), []byte("02"), []byte("03"), []byte("04")}

	for _, key := range keys[:2] {
		dup, err := nv.IsDuplicated(key)
		if err != nil {
			t.Error("failed check:", err)
			return
		}
		if dup {
			t.Error("key", string(key), "should not be marked as duplicated")
			return
		}
		err = nv.Save(key)
		if err != nil {
			t.Error("failed keep key:", err)
			return
		}
	}

	for _, key := range keys[:2] {
		dup, err := nv.IsDuplicated(key)
		if err != nil {
			t.Error("failed check:", err)
			return
		}
		if !dup {
			t.Error("key", string(key), "should be marked as duplicated")
			return
		}
	}

	// check cleanup
	// 1st: fill till factor (4 keys)
	for _, key := range keys[2:] {
		dup, err := nv.IsDuplicated(key)
		if err != nil {
			t.Error("failed check:", err)
			return
		}
		if dup {
			t.Error("key", string(key), "should not be marked as duplicated")
			return
		}
		err = nv.Save(key)
		if err != nil {
			t.Error("failed keep key:", err)
			return
		}
	}
	// 2nd: check amount of all keys
	allKeys, err := storages.AllKeysString(mem)
	if err != nil {
		t.Error("failed get all keys:", err)
		return
	}
	N := len(allKeys)
	if N != 2 {
		t.Error("not all keys removed: left", len(allKeys), allKeys)
		return
	}
	// 3rd: old one should now be non-dups

	for _, key := range keys[:2] {
		dup, err := nv.IsDuplicated(key)
		if err != nil {
			t.Error("failed check:", err)
			return
		}
		if dup {
			t.Error("key", string(key), "should not be marked as duplicated")
			return
		}
		err = nv.Save(key)
		if err != nil {
			t.Error("failed keep key:", err)
			return
		}
	}
}
