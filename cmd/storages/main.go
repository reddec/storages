package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/juju/fslock"
	"github.com/pkg/errors"
	"github.com/reddec/storages"
	"github.com/reddec/storages/cmd/storages/internal"
	storageconfig "github.com/reddec/storages/config"
	"github.com/reddec/storages/std"
	_ "github.com/reddec/storages/std/awsstorage"
	_ "github.com/reddec/storages/std/boltdb"
	_ "github.com/reddec/storages/std/filestorage"
	_ "github.com/reddec/storages/std/leveldbstorage"
	_ "github.com/reddec/storages/std/memstorage"
	_ "github.com/reddec/storages/std/redistorage"
	"github.com/reddec/storages/std/rest"
	_ "github.com/reddec/storages/std/rest"
	stor_utils "github.com/reddec/storages/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"time"
)

var version string = "dev"

type Config struct {
	URL       string        `short:"u" long:"url" env:"URL" description:"Storage URL" default:"bbolt://data"`
	Key       string        `short:"k" long:"key" env:"KEY" description:"Key in storage where configuration defined"`
	Lock      string        `short:"L" long:"lock" env:"LOCK" description:"Optional lock file for inter-process synchronization"`
	Supported listSupported `command:"supported" description:"list supported storages backends"`
	List      listKeys      `command:"list" alias:"ls" description:"list keys in storage"`
	Get       getKey        `command:"get" alias:"fetch" alias:"g" description:"get value by key"`
	Set       setKey        `command:"set" alias:"put" alias:"s" description:"set value for key"`
	Del       removeKey     `command:"remove" alias:"delete" alias:"del" alias:"rm" description:"remove value by key"`
	Copy      cpKeys        `command:"copy" alias:"cp" alias:"c" description:"copy keys from storage to destination"`
	Serve     restServe     `command:"serve" alias:"rest" description:"expose storage over REST interface"`
	Config    configCmd     `command:"config" alias:"cfg" description:"operations on configuration"`
	Queue     queueCmd      `command:"queue" alias:"q" description:"access to storage by naive queue interface"`
}

func (cfg *Config) getSource() storages.Storage {
	var lock *fslock.Lock
	if cfg.Lock != "" {
		lock = fslock.New(cfg.Lock)
		err := lock.Lock()
		if err != nil {
			log.Fatal("failed lock:", err)
		}
	}

	db, err := std.Create(config.URL)
	if err != nil {
		log.Fatal("failed initialize db:", err)
	}
	if lock != nil {
		db = stor_utils.WithCloseHook(db, func() {
			lock.Unlock()
		})
	}
	return db
}

func (cfg *Config) Storage() storages.Storage {
	src := cfg.getSource()
	if cfg.Key == "" {
		return src
	}
	stor, err := storageconfig.ParseJSON([]byte(cfg.Key), src)
	src.Close()
	if err != nil {
		log.Fatal("failed initialize db based on configuration from ", cfg.Key, ":", err)
	}
	return stor
}

var config Config

func main() {
	log.SetOutput(os.Stderr)
	parser := flags.NewParser(&config, flags.Default)
	parser.LongDescription = "Tools to work with storages\nAuthor: Baryshnikov Aleksandr <dev@baryshnikov.net>\nVersion: " + version
	_, err := parser.Parse()
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

type LineCodec interface {
	io.Closer
	Write([]byte) error
}

type listKeys struct {
	Null bool   `long:"null" short:"0" env:"NULL" description:"Use zero byte as terminator for list instead of new line (shorthand for -t null)"`
	JSON bool   `long:"json" env:"JSON" description:"Print keys as JSON array (shorthand for -t json)"`
	Type string `short:"t" long:"type" env:"TYPE" description:"Output encoding type" default:"plain" choice:"plain" choice:"null" choice:"json" choice:"base64" choice:"b64"`
}

func (l listKeys) getCodec() (LineCodec, error) {
	if l.Null {
		l.Type = "null"
	}
	if l.JSON {
		l.Type = "json"
	}
	switch l.Type {
	case "json":
		return internal.NewStringJSONLine(os.Stdout)
	case "null":
		return internal.NewPlainLine(os.Stdout, 0, false), nil
	case "plain":
		return internal.NewPlainLine(os.Stdout, '\n', true), nil
	case "base64", "b64":
		return internal.NewBase64Line(os.Stdout), nil
	default:
		return nil, errors.Errorf("encoding %v not known", l.Type)
	}
}

func (l listKeys) Execute(args []string) error {
	db := config.Storage()
	defer db.Close()
	codec, err := l.getCodec()
	if err != nil {
		return err
	}
	defer codec.Close()
	return db.Keys(codec.Write)
}

type getKey struct {
	Args struct {
		Key []string `description:"key names, if not set - STDIN lines used" positional-arg-name:"keys"`
	} `positional-args:"yes"`
}

func (g *getKey) Execute(args []string) error {
	db := config.Storage()
	defer db.Close()
	keys := getArgs(g.Args.Key...)
	for keys.Scan() {
		v, err := db.Get(keys.Bytes())
		if err != nil {
			return err
		}
		_, err = os.Stdout.Write(v)
		if err != nil {
			return err
		}
	}
	return os.Stdout.Close()
}

type setKey struct {
	Separator string `short:"s" long:"separator" env:"SEPARATOR" description:"Separator between key and value in line when stream used as source" default:" "`
	Args      struct {
		Key   string `description:"key name" positional-arg-name:"key"`
		Value string `description:"Value to put. Used STDIN if not set"`
	} `positional-args:"yes"`
}

func (s *setKey) Execute(args []string) error {
	db := config.Storage()
	defer db.Close()
	if len(s.Args.Key) > 0 {
		var data = []byte(s.Args.Value)
		if len(data) == 0 {
			v, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			data = v
		}
		return db.Put([]byte(s.Args.Key), data)
	} else {
		s.Separator = strings.ReplaceAll(s.Separator, "\\t", "\t")
		reader := bufio.NewScanner(os.Stdin)
		for reader.Scan() {
			line := reader.Bytes()
			if len(line) == 0 {
				continue
			}
			kv := bytes.SplitN(line, []byte(s.Separator), 2)
			if len(kv) == 1 {
				continue
			}
			err := db.Put(kv[0], kv[1])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type removeKey struct {
	Args struct {
		Key []string `description:"key names, if not set - STDIN lines used" positional-arg-name:"key"`
	} `positional-args:"yes"`
}

func (r *removeKey) Execute(args []string) error {
	db := config.Storage()
	defer db.Close()
	keys := getArgs(r.Args.Key...)
	for keys.Scan() {
		err := db.Del(keys.Bytes())
		if err != nil {
			return err
		}
	}
	return nil
}

type cpKeys struct {
	Args struct {
		URL string `description:"destination storage URL" positional-arg-name:"url" required:"yes"`
	} `positional-args:"yes"`
}

func (c *cpKeys) Execute(args []string) error {
	from := config.Storage()
	defer from.Close()
	to, err := std.Create(c.Args.URL)
	if err != nil {
		return err
	}
	defer to.Close()
	return from.Keys(func(key []byte) error {
		data, err := from.Get(key)
		if err != nil {
			return err
		}
		return to.Put(key, data)
	})
}

type restServe struct {
	GracefulShutdown time.Duration `long:"graceful-shutdown" env:"GRACEFUL_SHUTDOWN" description:"Interval before server shutdown" default:"15s"`
	Bind             string        `long:"bind" env:"BIND" description:"Address to where bind HTTP server" default:"0.0.0.0:8080"`
	TLS              bool          `long:"tls" env:"TLS" description:"Enable HTTPS serving with TLS"`
	CertFile         string        `long:"cert-file" env:"CERT_FILE" description:"Path to certificate for TLS" default:"server.crt"`
	KeyFile          string        `long:"key-file" env:"KEY_FILE" description:"Path to private key for TLS" default:"server.key"`
}

func (r *restServe) Execute(args []string) error {
	storage := config.Storage()
	defer storage.Close()

	server := http.Server{
		Addr:    r.Bind,
		Handler: rest.NewServer(storage),
	}

	go func() {
		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Kill, os.Interrupt)
		<-c
		ctx, cancel := context.WithTimeout(context.Background(), r.GracefulShutdown)
		defer cancel()
		server.Shutdown(ctx)
	}()
	log.Println("REST server is on", r.Bind)
	if r.TLS {
		return server.ListenAndServeTLS(r.CertFile, r.KeyFile)
	}
	return server.ListenAndServe()
}

func getArgs(def ...string) *bufio.Scanner {
	if len(def) == 0 {
		return bufio.NewScanner(os.Stdin)
	}
	return bufio.NewScanner(bytes.NewBufferString(strings.Join(def, "\n")))
}
