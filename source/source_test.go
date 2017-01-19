package source

import (
	"reflect"
	"testing"
)

func TestSource(t *testing.T) {
	name := "sourceName"
	tag := "tag"
	url := "http://url.com"

	nts := NewSource(name, tag, url)

	t.Run("Source has correct name, tag, endpoint", func(t *testing.T) {
		if nts.name != name || nts.tag != tag || nts.Endpoint != url {
			t.Errorf("Expected %s, %s, %s, got %v", name, tag, url, nts)
		}
	})

	t.Run("Source has default client", func(t *testing.T) {
		defaultClient := defaultClient()
		if !reflect.DeepEqual(defaultClient, nts.Client) {
			t.Errorf("Expected %v, got: %v", defaultClient, nts.Client)
		}
	})

	username := "test"
	password := "test1"

	nts.SetAuth(username, password)
	t.Run("Source has correct Username and Password", func(t *testing.T) {
		if nts.Username != username || nts.Password != password {
			t.Errorf("Expected %s, %s, got %s, %s",
				username, password,
				nts.Username, nts.Password)
		}
	})

	t.Run("Source has correct SourceOptions", func(t *testing.T) {

	})
}
