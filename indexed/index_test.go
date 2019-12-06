package indexed

import (
	"encoding/json"
	"fmt"
	"github.com/reddec/storages/std/memstorage"
)

func ExampleNewIndex() {
	type User struct {
		ID    string // primary key
		Email string // unique
		Year  string // non-unique, but indexed
	}
	alice := User{
		ID:    "1",
		Email: "alice@example.com",
		Year:  "1935",
	}

	bob := User{
		ID:    "2",
		Email: "bob@example.com",
		Year:  "1935",
	}

	// create primary index (storage)
	primary := memstorage.New()
	byEmail := NewUniqueIndex(memstorage.New())
	byYear := NewIndex(memstorage.New())

	// save data, indexed by ID
	_ = primary.Put([]byte(alice.ID), toJSON(alice))
	_ = primary.Put([]byte(bob.ID), toJSON(bob))

	// add refs to indexes
	_ = byEmail.Link([]byte(alice.ID), []byte(alice.Email))
	_ = byEmail.Link([]byte(bob.ID), []byte(bob.Email))

	_ = byYear.Link([]byte(alice.ID), []byte(alice.Year))
	_ = byYear.Link([]byte(bob.ID), []byte(bob.Year))

	// find alice by email
	pKeys, _ := byEmail.Find([]byte(alice.Email))
	userData, _ := primary.Get(pKeys[0])

	var user User
	fromJSON(userData, &user)
	fmt.Println("records:", len(pKeys))
	fmt.Println("user id:", user.ID)
	fmt.Println("")

	// find all users with year 1935
	pKeys, _ = byYear.Find([]byte("1935"))
	fmt.Println("records:", len(pKeys))
	for _, pk := range pKeys {
		userData, _ = primary.Get(pk)
		var user User
		fromJSON(userData, &user)
		fmt.Println("user id:", user.ID)
	}
	// Output:
	// records: 1
	// user id: 1
	//
	// records: 2
	// user id: 1
	// user id: 2
}

func toJSON(data interface{}) []byte {
	d, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return d
}

func fromJSON(data []byte, item interface{}) {
	err := json.Unmarshal(data, item)
	if err != nil {
		panic(err)
	}
}
