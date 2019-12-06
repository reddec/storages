package std

import (
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"github.com/reddec/storages"
	"net/url"
	"sync"
)

// Plain configuration for storage. Exported field will mapped automatically
type Configuration interface {
	// Create new instance of storage or fail
	Create() (storages.Storage, error)
}

// Factory function for configuration
type FactoryFunc func() Configuration

// Factory function with custom mapping
type FactoryURLFunc func(*url.URL) (storages.Storage, error)

var supported sync.Map

// Register new factory for defined schema with github.com/gorilla/schema mapper for query parameters
func Register(schema string, factoryFunc FactoryFunc) {
	supported.Store(schema, factoryFunc)
}

// Register new factory for defined schema with custom URL mapping logic
func RegisterWithMapper(schema string, factoryFunc FactoryURLFunc) {
	supported.Store(schema, factoryFunc)
}

// Create new storage by looking into storages registry and mapping url parameters to configuration
func Create(rawURL string) (storages.Storage, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	factory, ok := supported.Load(u.Scheme)
	if !ok {
		return nil, errors.Errorf("unsupported storage scheme: %v", u.Scheme)
	}
	if customFactory, ok := factory.(FactoryURLFunc); ok {
		return customFactory(u)
	}
	configTemplate := factory.(FactoryFunc)()
	err = schema.NewDecoder().Decode(configTemplate, u.Query())
	if err != nil {
		return nil, err
	}
	return configTemplate.Create()
}

// Supported schemas that depends of imports.
//
// Use import like `_ "github.com/reddec/storages/std/rest"`
func Supported() []string {
	var ans []string
	supported.Range(func(key, value interface{}) bool {
		ans = append(ans, key.(string))
		return true
	})
	return ans
}
