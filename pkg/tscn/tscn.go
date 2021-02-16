// Package tscn parses and returns a representation of Godots TSCN file format
package tscn

import (
	"io"

	"github.com/pkg/errors"

	"github.com/atomicptr/godot-tscn-parser/internal/convert"
	"github.com/atomicptr/godot-tscn-parser/internal/parser"
	"github.com/atomicptr/godot-tscn-parser/internal/validate"
	"github.com/atomicptr/godot-tscn-parser/pkg/godot"
)

// ParseScene parses a TSCN file of the type gd_scene
func ParseScene(r io.Reader) (*godot.Scene, error) {
	tscn, err := parseAndValidateTscnFile(r)
	if err != nil {
		return nil, err
	}
	return convert.ToGodotScene(tscn)
}

// ParseProject parses the central project.godot project configuration file
func ParseProject(r io.Reader) (*godot.Project, error) {
	tscn, err := parseAndValidateTscnFile(r)
	if err != nil {
		return nil, err
	}
	return convert.ToGodotProject(tscn)
}

// ParseResource parses .tres files
func ParseResource(r io.Reader) (*godot.Resource, error) {
	tscn, err := parseAndValidateTscnFile(r)
	if err != nil {
		return nil, err
	}
	return convert.ToGodotResource(tscn)
}

func parseAndValidateTscnFile(r io.Reader) (*parser.TscnFile, error) {
	tscn, err := parser.Parse(r)
	if err != nil {
		return nil, errors.Wrap(err, "parser error")
	}
	if err = validate.TscnFileFormat(tscn); err != nil {
		return nil, errors.Wrap(err, "invalid file format")
	}
	return tscn, nil
}
