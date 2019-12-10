package filestorage

import (
	"encoding/json"
	"github.com/reddec/storages"
	"github.com/reddec/storages/std"
	"net/url"
	"os"
	"path/filepath"
	"sync"
)

// Converter from some data to bytes
type EncoderFunc func(value interface{}) ([]byte, error)

// Decode bytes to value
type DecoderFunc func(data []byte, value interface{}) error

// New single file storage with custom encoder and decoder
func NewEncodedFile(filename string, encoderFunc EncoderFunc, decoderFunc DecoderFunc) *encodedNamespace {
	stor := &encodedStorage{
		encoder:  encoderFunc,
		decoder:  decoderFunc,
		filename: filename,
	}
	stor.root = encodedNamespace{
		data:    &dataType{},
		storage: stor,
	}
	return &stor.root
}

// New single file storage with JSON encoding
func NewJSONFile(filename string) *encodedNamespace {
	return NewEncodedFile(filename, func(value interface{}) (bytes []byte, err error) {
		return json.MarshalIndent(value, "", "  ")
	}, json.Unmarshal)
}

type dataType struct {
	Data       map[string][]byte    `json:"data" xml:"data" yaml:"data"`
	Namespaces map[string]*dataType `json:"namespaces,omitempty" xml:"namespaces,omitempty" yaml:"namespaces,omitempty"`
}

type encodedStorage struct {
	lock     sync.RWMutex
	encoder  EncoderFunc
	decoder  DecoderFunc
	filename string
	root     encodedNamespace
}

type encodedNamespace struct {
	lock    sync.RWMutex
	data    *dataType
	storage *encodedStorage
}

func (e *encodedNamespace) Put(key []byte, data []byte) error {
	cp := make([]byte, len(data))
	copy(cp, data)
	e.lock.Lock()
	defer e.lock.Unlock()

	if e.data.Data == nil {
		e.data.Data = make(map[string][]byte)
	}
	e.data.Data[string(key)] = cp
	return e.storage.safeDump()
}

func (e *encodedNamespace) Close() error { return nil }

func (e *encodedNamespace) Get(key []byte) ([]byte, error) {
	e.lock.RLock()
	defer e.lock.RUnlock()
	v, ok := e.data.Data[string(key)]
	if !ok {
		return nil, os.ErrNotExist
	}
	cp := make([]byte, len(v))
	copy(cp, v)
	return cp, nil
}

func (e *encodedNamespace) Del(key []byte) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	k := string(key)
	if _, ok := e.data.Data[k]; !ok {
		return nil
	}
	delete(e.data.Data, k)
	return e.storage.safeDump()
}

func (e *encodedNamespace) Keys(handler func(key []byte) error) error {
	e.lock.RLock()
	defer e.lock.RUnlock()
	for k := range e.data.Data {
		err := handler([]byte(k))
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *encodedNamespace) Namespace(name []byte) (storages.Storage, error) {
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.data.Namespaces == nil {
		e.data.Namespaces = make(map[string]*dataType)
	}
	ns, ok := e.data.Namespaces[string(name)]
	if !ok {
		ns = &dataType{}
		e.data.Namespaces[string(name)] = ns
	}
	return &encodedNamespace{
		data:    ns,
		storage: e.storage,
	}, nil
}

func (e *encodedNamespace) Namespaces(handler func(name []byte) error) error {
	e.lock.RLock()
	defer e.lock.RUnlock()
	for k := range e.data.Namespaces {
		err := handler([]byte(k))
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *encodedNamespace) DelNamespace(name []byte) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	k := string(name)
	_, ok := e.data.Namespaces[k]
	if !ok {
		return nil
	}
	delete(e.data.Namespaces, k)
	return e.storage.safeDump()
}

func (e *encodedStorage) safeDump() error {
	e.lock.Lock()
	defer e.lock.Unlock()
	bin, err := e.encoder(e.root.data)
	if err != nil {
		return err
	}
	return safeWrite(e.filename, bin)
}

func init() {
	std.RegisterWithMapper("file+json", func(url *url.URL) (storage storages.Storage, e error) {
		return NewJSONFile(filepath.Join(url.Host, url.Path)), nil
	})
}
