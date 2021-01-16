package tscn

import (
	"github.com/atomicptr/godot-tscn-parser/pkg/parser"
	"io"
	"os"
)

// TODO: do not use parser.GdScene here... build a proper node tree
func Parse(r io.Reader) (*parser.TscnFile, error) {
	return parser.Parse(r)
}

// TODO: do not use parser.GdScene here... build a proper node tree
func LoadFileAndParse(file string) (*parser.TscnFile, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	return Parse(f)
}
