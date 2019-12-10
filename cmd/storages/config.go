package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	config2 "github.com/reddec/storages/config"
	"os"
)

type configCmd struct {
	Example configExample `command:"example" description:"generate examples for different kind of storages and save it to current storage"`
}

type configExample struct {
	Type string `short:"t" long:"type" env:"TYPE" description:"Storage type for example" default:"redundant" choice:"simple" choice:"redundant" choice:"sharded"`
}

func (c *configExample) Execute(args []string) error {
	var data []byte
	var err error
	switch c.Type {
	case "simple":
		data, err = json.MarshalIndent(config2.Simple{URL: "bbolt://data1"}, "", " ")
	case "sharded":
		data, err = json.MarshalIndent(config2.Sharded{
			Shards: []string{
				"key1forConfigurationForStorage",
				"key2forConfigurationForStorage",
			},
		}, "", " ")
	case "redundant":
		rdr := config2.Redundant{
			Read: config2.ReadStrategy{
				First: &struct{}{},
			},
			Write: config2.WriteStrategy{
				AtLeast: &config2.AtLeast{Num: 2},
			},
			Dedup: "key1forStorageForDeduplication",
			Storages: []string{
				"key2forConfigurationForStorage",
				"key3forConfigurationForStorage",
			},
		}
		data, err = json.MarshalIndent(rdr, "", "  ")
	default:
		return errors.Errorf("unknown type %v", c.Type)
	}
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(data)
	return err
}
