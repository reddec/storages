package tests

import (
	"github.com/reddec/storages/std"
	_ "github.com/reddec/storages/std/redistorage"
	_ "github.com/reddec/storages/std/rest"
	"reflect"
	"testing"
)

func TestCreateByURL(t *testing.T) {
	testCreate("http://example.com/abc", "restClient", t)
	testCreate("redis://myhost/1?key=data", "redisStorage", t)
}

func testCreate(url string, expected string, t *testing.T) {
	instance, err := std.Create(url)
	if err != nil {
		t.Error(expected+" from "+url+":", err)
		return
	}
	if tp := typename(instance); tp != expected {
		t.Errorf(expected+" from "+url+" is not a "+expected+": %v", tp)
	}
}

func typename(v interface{}) string {
	return reflect.ValueOf(v).Elem().Type().Name()
}
