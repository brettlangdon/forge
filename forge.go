package forge

import (
	"bytes"
	"io"
	"strings"

	"github.com/brettlangdon/forge/config"
	"github.com/brettlangdon/forge/parser"
)

func ParseString(data string) (*config.SectionValue, error) {
	return parser.ParseReader(strings.NewReader(data))
}

func ParseBytes(data []byte) (*config.SectionValue, error) {
	return parser.ParseReader(bytes.NewReader(data))
}

func ParseFile(filename string) (*config.SectionValue, error) {
	return parser.ParseFile(filename)
}

func ParseReader(reader io.Reader) (*config.SectionValue, error) {
	return parser.ParseReader(reader)
}
