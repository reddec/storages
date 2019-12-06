package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/reddec/storages"
	"github.com/reddec/storages/std"
	_ "github.com/reddec/storages/std/awsstorage"
	_ "github.com/reddec/storages/std/boltdb"
	_ "github.com/reddec/storages/std/filestorage"
	_ "github.com/reddec/storages/std/leveldbstorage"
	_ "github.com/reddec/storages/std/memstorage"
	_ "github.com/reddec/storages/std/redistorage"
	_ "github.com/reddec/storages/std/rest"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

type Config struct {
	URL       string        `short:"u" long:"url" env:"URL" description:"Storage URL" default:"bbolt://data"`
	Supported listSupported `command:"supported" description:"list supported storages backends"`
	List      listKeys      `command:"list" alias:"ls" description:"list keys in storage"`
	Get       getKey        `command:"get" alias:"fetch" alias:"g" description:"get value by key"`
	Set       setKey        `command:"set" alias:"put" alias:"s" description:"set value for key"`
	Del       removeKey     `command:"remove" alias:"delete" alias:"del" alias:"rm" description:"remove value by key"`
}

func (cfg *Config) Storage() storages.Storage {
	db, err := std.Create(config.URL)
	if err != nil {
		log.Fatal("failed initialize db:", err)
	}
	return db
}

var config Config

func main() {
	_, err := flags.Parse(&config)
	if err != nil {
		os.Exit(1)
	}
}

type listSupported struct{}

func (l listSupported) Execute(args []string) error {
	names := std.Supported()
	sort.Strings(names)
	for _, name := range names {
		fmt.Println(name)
	}
	return nil
}

type listKeys struct {
	Null bool `long:"null" short:"0" env:"NULL" description:"Use zero byte as terminator for list instead of new line"`
}

func (l listKeys) Execute(args []string) error {
	db := config.Storage()
	defer db.Close()
	return db.Keys(func(key []byte) error {
		var err error
		_, err = os.Stdout.Write(key)
		if err != nil {
			return err
		}
		if l.Null {
			_, err = os.Stdout.Write([]byte{0})
		} else {
			_, err = os.Stdout.Write([]byte("\n"))
		}
		return err
	})
}

type getKey struct {
	Args struct {
		Key string `description:"key name" positional-arg-name:"key" required:"yes"`
	} `positional-args:"yes"`
}

func (g *getKey) Execute(args []string) error {
	db := config.Storage()
	defer db.Close()
	v, err := db.Get([]byte(g.Args.Key))
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(v)
	return err
}

type setKey struct {
	Stream bool `long:"stream" short:"s" env:"STREAM" description:"Use STDIN as source of value"`
	Args   struct {
		Key   string `description:"key name" positional-arg-name:"key" required:"yes"`
		Value string `description:"Value to put if stream flag is not enabled"`
	} `positional-args:"yes"`
}

func (s *setKey) Execute(args []string) error {
	var data = []byte(s.Args.Value)
	if s.Stream || len(data) == 0 {
		v, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		data = v
	}
	db := config.Storage()
	defer db.Close()
	return db.Put([]byte(s.Args.Key), data)
}

type removeKey struct {
	Args struct {
		Key string `description:"key name" positional-arg-name:"key" required:"yes"`
	} `positional-args:"yes"`
}

func (r *removeKey) Execute(args []string) error {
	db := config.Storage()
	defer db.Close()
	return db.Del([]byte(r.Args.Key))
}
