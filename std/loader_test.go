package std

import _ "github.com/reddec/storages/std/redistorage"

func ExampleCreate() {
	storage, err := Create("redis://my-host")
	if err != nil {
		panic(err)
	}
	defer storage.Close()
}
