package validate

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/atomicptr/godot-tscn-parser/internal/parser"
)

// https://docs.godotengine.org/en/stable/development/file_formats/tscn.html#external-resources
func validatorExtResourceRequiredAttributes(tscnFile *parser.TscnFile) error {
	for _, section := range tscnFile.Sections {
		if section.ResourceType != parser.ResourceTypeExtResource {
			continue
		}

		// validate path
		path, err := section.GetAttribute("path")
		if err != nil {
			return errors.Wrapf(err, "ext_resource is missing required field 'path' %s", section.Pos)
		}
		if path.String == nil {
			return fmt.Errorf("ext_resource attribute path must be a string %s", section.Pos)
		}

		// validate type
		t, err := section.GetAttribute("type")
		if err != nil {
			return errors.Wrapf(err, "ext_resource is missing required field 'type' %s", section.Pos)
		}
		if t.String == nil {
			return fmt.Errorf("ext_resource attribute type must be a string %s", section.Pos)
		}

		// validate id
		id, err := section.GetAttribute("id")
		if err != nil {
			return errors.Wrapf(err, "ext_resource is missing required field 'id' %s", section.Pos)
		}
		if id.Integer == nil {
			return fmt.Errorf("ext_resource attribute id must be an integer %s", section.Pos)
		}
	}

	return nil
}
