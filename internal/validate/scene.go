package validate

import (
	"fmt"

	"github.com/atomicptr/godot-tscn-parser/internal/parser"
	"github.com/atomicptr/godot-tscn-parser/pkg/godot"
)

func validatorSceneIsInSupportedFormat(tscnFile *parser.TscnFile) error {
	version, err := tscnFile.GetAttribute("format")
	if err != nil {
		return nil // no format specified, don't bother
	}

	if version.Integer == nil {
		return fmt.Errorf("gd_scene format is not an integer %s", tscnFile.Pos)
	}

	if *version.Integer != godot.FormatVersion {
		return fmt.Errorf(
			"gd_scene format is unsupported version '%d', we only support version '%d' %s",
			*version.Integer,
			godot.FormatVersion,
			tscnFile.Pos,
		)
	}

	return nil
}
