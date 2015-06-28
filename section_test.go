package forge_test

import (
	"testing"

	"github.com/brettlangdon/forge"
)

func TestSectionKeys(t *testing.T) {
	t.Parallel()

	section := forge.NewSection()
	section.SetString("key1", "value1")
	section.SetString("key2", "value2")
	section.SetString("key3", "value3")

	keys := section.Keys()

	if len(keys) != 3 {
		t.Error("expected Section to have 3 keys")
	}

	if keys[0] != "key1" {
		t.Error("expected 'key1' to be in the list of keys")
	}
	if keys[1] != "key2" {
		t.Error("expected 'key2' to be in the list of keys")
	}
	if keys[2] != "key3" {
		t.Error("expected 'key3' to be in the list of keys")
	}
}
