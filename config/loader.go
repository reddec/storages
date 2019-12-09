package config

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/reddec/storages"
	"github.com/reddec/storages/dedup"
	"github.com/reddec/storages/sharded"
	"github.com/reddec/storages/std"
	"github.com/reddec/storages/std/memstorage"
)

// Decoding function for saved configuration
type DecoderFunc func(in []byte, out interface{}) error

// Recursively parse configuration from start (main) configuration defined in entry cell. With standard JSON decoder
func ParseJSON(entryKey []byte, storage storages.Storage) (storages.Storage, error) {
	return ParseWithDecoder(entryKey, storage, json.Unmarshal)
}

// Recursively parse configuration from start (main) configuration defined in entry cell. With custom decoder
func ParseWithDecoder(entryKey []byte, storage storages.Storage, decoderFunc DecoderFunc) (storages.Storage, error) {
	var loadedStorages = map[string]storages.Storage{}
	return getStorage(entryKey, storage, decoderFunc, loadedStorages)
}

func getStorage(key []byte, storage storages.Storage, decoderFunc DecoderFunc, loaded map[string]storages.Storage) (storages.Storage, error) {
	if saved, ok := loaded[string(key)]; ok {
		return saved, nil
	}
	data, err := storage.Get(key)
	if err != nil {
		return nil, err
	}
	var kind Kind
	err = decoderFunc(data, &kind)
	if err != nil {
		return nil, err
	}

	switch kind.Kind {
	case "simple":
		var simple Simple
		err = decoderFunc(data, &simple)
		if err != nil {
			return nil, err
		}
		stor, err := std.Create(simple.URL)
		if err != nil {
			return nil, err
		}
		loaded[string(key)] = stor
		return stor, nil
	case "sharded":
		var cfg Sharded
		err = decoderFunc(data, &cfg)
		if err != nil {
			return nil, err
		}
		var shards []storages.Storage
		for _, shardId := range cfg.Shards {
			shard, err := getStorage([]byte(shardId), storage, decoderFunc, loaded)
			if err != nil {
				return nil, errors.Wrapf(err, "get shard %v", shardId)
			}
			shards = append(shards, shard)
		}
		stor := storages.Sharded(sharded.NewHashedArray(shards))
		loaded[string(key)] = stor
		return stor, nil
	case "redundant":
		var redundant Redundant
		err = decoderFunc(data, &redundant)
		if err != nil {
			return nil, err
		}
		var backs []storages.Storage
		for _, storageID := range redundant.Storages {
			back, err := getStorage([]byte(storageID), storage, decoderFunc, loaded)
			if err != nil {
				return nil, errors.Wrapf(err, "get backed storage %v", storageID)
			}
			backs = append(backs, back)
		}
		var dedupStorage storages.Storage
		if redundant.Dedup != "" {
			dedupStorage, err = getStorage([]byte(redundant.Dedup), storage, decoderFunc, loaded)
			if err != nil {
				return nil, errors.Wrapf(err, "get deduplication storage %v", redundant.Dedup)
			}
		} else {
			dedupStorage = memstorage.New()
		}
		return storages.Redundant(redundant.Write.GetStrategy(backs), redundant.Read.GetStrategy(backs), dedup.Offloaded(dedupStorage), backs...), nil
	default:
		return nil, errors.Errorf("unknown storage kind '%v' in %v", kind.Kind, string(key))
	}

}
