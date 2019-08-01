package stor_utils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"reddec/storages"
)

// Put data to storage using default JSON encoder
func PutJSON(st storages.Writer, key string, data interface{}) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return st.Put([]byte(key), d)
}

// Get data from storage and (if persist) decode it to target using default JSON decoder
func GetJSON(st storages.Storage, key string, target interface{}) error {
	val, err := st.Get([]byte(key))
	if err != nil {
		return err
	}
	return json.Unmarshal(val, target)
}

// Put data to storage using GOB encoder
func PutGOB(st storages.Writer, key string, data interface{}) error {
	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(data)
	if err != nil {
		return err
	}
	return st.Put([]byte(key), buf.Bytes())
}

// Get data from storage and (if persist) decode it to target using GOB decoder
func GetGOB(st storages.Storage, key string, target interface{}) error {
	val, err := st.Get([]byte(key))
	if err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(val)).Decode(target)
}
