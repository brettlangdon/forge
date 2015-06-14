package forge

import (
	"bytes"
	"io"
	"strings"

	"github.com/brettlangdon/forge/parser"
)

func ParseString(data string) (map[string]interface{}, error) {
	settings, err := parser.ParseReader(strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	return settings.ToMap()
}

func ParseBytes(data []byte) (map[string]interface{}, error) {
	settings, err := parser.ParseReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return settings.ToMap()
}

func ParseFile(filename string) (map[string]interface{}, error) {
	settings, err := parser.ParseFile(filename)
	if err != nil {
		return nil, err
	}

	return settings.ToMap()
}

func ParseReader(reader io.Reader) (map[string]interface{}, error) {
	settings, err := parser.ParseReader(reader)
	if err != nil {
		return nil, err
	}

	return settings.ToMap()
}
