package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"github.com/reddec/storages"
	"github.com/reddec/storages/awsstorage"
	"github.com/reddec/storages/filestorage"
	"github.com/reddec/storages/leveldbstorage"
	"github.com/reddec/storages/redistorage"
	"io/ioutil"
	"log"
	"os"
)

type source interface {
	Build() (storages.Storage, error)
}

type fileParams struct {
	Location string `long:"location" env:"LOCATION" description:"Root dir to store data" default:"./db"`
	Flat     bool   `long:"flat" env:"FLAT" description:"Use flat file organization: key name is equal to file (prone to invalid key names)"`
}

func (fp *fileParams) Build() (storages.Storage, error) {
	if fp.Flat {
		return filestorage.NewFlat(fp.Location), nil
	}
	return filestorage.NewDefault(fp.Location), nil
}

type levelDbParams struct {
	Location string `long:"location" env:"LOCATION" description:"Root dir to store data" default:"./db"`
}

func (ld *levelDbParams) Build() (storages.Storage, error) {
	return leveldbstorage.New(ld.Location)
}

type redisParams struct {
	URL       string `long:"url" env:"URL" description:"Redis URL" default:"redis://localhost"`
	Namespace string `long:"namespace" env:"NAMESPACE" description:"Hashmap name" default:"db"`
}

func (rd *redisParams) Build() (storages.Storage, error) {
	return redistorage.New(rd.Namespace, rd.URL)
}

type awsParams struct {
	Bucket         string `long:"bucket" env:"BUCKET" description:"S3 AWS bucket"`
	Endpoint       string `long:"endpoint" env:"ENDPOINT" description:"Override AWS endpoint for AWS-capable services"`
	ForcePathStyle bool   `long:"force-path-style" env:"FORCE_PATH_STYLE" description:"Force the request to use path-style addressing"`
}

func (ap *awsParams) Build() (storages.Storage, error) {
	config := aws.NewConfig()
	if ap.Endpoint != "" {
		config = config.WithEndpoint(ap.Endpoint)
	}
	config.S3ForcePathStyle = &ap.ForcePathStyle
	config.Credentials = credentials.NewEnvCredentials()
	return awsstorage.New(ap.Bucket, config)
}

var config struct {
	Db      string        `long:"db" short:"t" env:"DB" description:"DB mode" default:"file" choice:"file" choice:"leveldb" choice:"redis" choice:"s3"`
	Stream  bool          `long:"stream" short:"s" env:"STREAM" description:"Use STDIN as source of value"`
	Null    bool          `long:"null" short:"0" env:"NULL" description:"Use zero byte as terminator for list instead of new line"`
	File    fileParams    `group:"File storage params" namespace:"file" env-namespace:"FILE"`
	LevelDB levelDbParams `group:"LevelDB storage params" namespace:"leveldb" env-namespace:"LEVELDB"`
	Redis   redisParams   `group:"Redis storage params" namespace:"redis" env-namespace:"REDIS"`
	S3      awsParams     `group:"S3 storage" namespace:"s3" env-namespace:"S3"`
	Args    struct {
		Command string `description:"what to do (put/set, list, get, del)" choice:"get" choice:"put" choice:"set" choice:"list" choice:"ls" choice:"del" default:"list" required:"yes"`
		Key     string `description:"key name" positional-arg-name:"key"`
		Value   string `description:"Value to put if stream flag is not enabled"`
	} `positional-args:"yes"`
}

func main() {
	_, err := flags.Parse(&config)
	if err != nil {
		os.Exit(1)
	}

	var src source
	switch config.Db {
	case "file":
		src = &config.File
	case "leveldb":
		src = &config.LevelDB
	case "redis":
		src = &config.Redis
	case "s3":
		src = &config.S3
	default:
		panic("unknown db type: " + config.Db)
	}

	db, err := src.Build()
	if err != nil {
		log.Fatal("failed initialize", config.Db, "db:", err)
	}

	err = run(db)
	db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func run(db storages.Storage) error {
	var data = []byte(config.Args.Value)
	if (config.Stream || len(data) == 0) && config.Args.Command == "put" {
		v, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		data = v
	}
	switch config.Args.Command {
	case "put", "set":
		return db.Put([]byte(config.Args.Key), data)
	case "get":
		v, err := db.Get([]byte(config.Args.Key))
		if err != nil {
			return err
		}
		os.Stdout.Write(v)
		return nil
	case "list", "ls":
		return db.Keys(func(key []byte) error {
			os.Stdout.Write(key)
			if config.Null {
				os.Stdout.Write([]byte{0})
			} else {
				os.Stdout.Write([]byte("\n"))
			}
			return nil
		})
	default:
		return errors.New("Unknown command " + config.Args.Command)
	}
}
