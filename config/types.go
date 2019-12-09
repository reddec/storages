package config

import (
	"github.com/reddec/storages"
)

// Kind of storage: sharded, simple or redundant
type Kind struct {
	Kind string `json:"kind" yaml:"kind" xml:"kind"` // storage kind name
}

// Simple storage config - a storage that could be initialized by URL
type Simple struct {
	URL string `json:"url" yaml:"url" xml:"url"` // storage URL for initialization (see std package)
}

// Sharded storage config
type Sharded struct {
	// names of underlying storages
	// that will be initialized and used as shards
	Shards []string `json:"shards" yaml:"shards" xml:"shards"`
}

// Redundant storage config
type Redundant struct {
	Read  ReadStrategy  `json:"read" yaml:"read" xml:"read"`    // how to read data
	Write WriteStrategy `json:"write" yaml:"write" xml:"write"` // how to write data
	// optional storage used for deduplication during keys iteration
	Dedup string `json:"dedup" yaml:"dedup" xml:"dedup"`
	// names of underlying storages
	// that will be initialized and used for distribution
	Storages []string `json:"storages" yaml:"storages" xml:"storages"`
}

// Read strategy config for redundant storage
type ReadStrategy struct {
	// Iterate over storages until first value returned without error
	First *struct {
	} `json:"first" yaml:"first" xml:"first"`
}

// Initialize strategy. If no strategy defined - used `First` strategy
func (rrs ReadStrategy) GetStrategy(backs []storages.Storage) storages.DReader {
	switch {
	case rrs.First != nil:
		return storages.First()
	default:
		return storages.First()
	}
}

// Write strategy config for redundant storage
type WriteStrategy struct {
	// At least specified amount of successful writes should be done for success
	AtLeast *struct {
		Num int `json:"num" yaml:"num" xml:"num"` // minimal amount of successful write
	} `json:"atleast" yaml:"atleast" xml:"atleast"`
}

// Initialize strategy. If no strategy defined - all writes should be complete without error for success
func (rrs WriteStrategy) GetStrategy(backs []storages.Storage) storages.DWriter {
	switch {
	case rrs.AtLeast != nil:
		return storages.AtLeast(rrs.AtLeast.Num)
	default:
		return storages.AtLeast(len(backs))
	}
}
