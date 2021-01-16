// Package tscn parses and returns a representation of Godots TSCN file format
package tscn

import (
	"github.com/atomicptr/godot-tscn-parser/pkg/parser"
	"io"
	"os"
	"path/filepath"
)

// Parse a TSCN file and returns a struct representing the files content
// TODO: do not use parser.GdScene here... build a proper node tree
func Parse(r io.Reader) (*parser.TscnFile, error) {
	return parser.Parse(r)
}

// LoadFileAndParse does what it says it does, check Parse for more information.
// TODO: do not use parser.GdScene here... build a proper node tree
func LoadFileAndParse(file string) (*parser.TscnFile, error) {
	f, err := os.Open(filepath.Clean(file))
	if err != nil {
		return nil, err
	}

	return Parse(f)
}
