// Package validate does TSCN validate
package validate

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/atomicptr/godot-tscn-parser/internal/parser"
)

type tscnFileValidationFunc func(tscnFile *parser.TscnFile) error

// TscnFileFormat runs a set of pre-defined validators against a TSCN file
func TscnFileFormat(tscnFile *parser.TscnFile) error {
	validators := []struct {
		Name     string
		Function tscnFileValidationFunc
	}{
		{"ExtResource has required attributes", validatorExtResourceRequiredAttributes},
		{"Scene root must not have a path attribute", validatorFirstNodeHasNoParent},
		{"Scene must not have multiple root nodes", validatorOnlyOneRootNode},
		{"Scene must be set to supported version", validatorSceneIsInSupportedFormat},
		{"All references to ExtResource/SubResource must exist", validatorTestIfAllResourceReferencesActuallyExist},
	}

	for _, validator := range validators {
		err := validator.Function(tscnFile)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("validate for '%s' failed", validator.Name))
		}
	}

	return nil
}
