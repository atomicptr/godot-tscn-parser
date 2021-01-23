// Package tscn parses and returns a representation of Godots TSCN file format
package tscn

import (
	"github.com/atomicptr/godot-tscn-parser/internal/parser"
	"github.com/atomicptr/godot-tscn-parser/pkg/godot"
	"github.com/pkg/errors"
	"io"
)

// ParseScene parses a TSCN file of the type gd_scene
func ParseScene(r io.Reader) (*godot.Scene, error) {
	tscn, err := parser.Parse(r)
	if err != nil {
		return nil, errors.Wrap(err, "parser error")
	}
	return tscn.ConvertToGodotScene()
}
