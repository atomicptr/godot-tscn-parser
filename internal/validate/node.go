package validate

import (
	"fmt"

	"github.com/atomicptr/godot-tscn-parser/internal/parser"
)

func validatorFirstNodeHasNoParent(tscnFile *parser.TscnFile) error {
	for _, section := range tscnFile.Sections {
		if section.ResourceType != parser.ResourceTypeNode {
			continue
		}

		_, err := section.GetAttribute("parent")
		if err == nil {
			return fmt.Errorf(
				"the first node in the file, which is also the scene root, must not have a 'parent' attribute. %s",
				section.Pos,
			)
		}
		// we only care about the first node
		return nil
	}

	return nil
}

func validatorOnlyOneRootNode(tscnFile *parser.TscnFile) error {
	foundRoot := false
	for _, section := range tscnFile.Sections {
		if section.ResourceType != parser.ResourceTypeNode {
			continue
		}

		_, err := section.GetAttribute("parent")
		if err != nil {
			if foundRoot {
				return fmt.Errorf("found a second root node (a node without parent): %s", section.Pos)
			}

			foundRoot = true
		}
	}

	return nil
}
