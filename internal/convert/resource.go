package convert

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/atomicptr/godot-tscn-parser/internal/parser"
	"github.com/atomicptr/godot-tscn-parser/pkg/godot"
)

// ToGodotResource tries to convert a TscnFile structure to a .tres file
func ToGodotResource(tscn *parser.TscnFile) (*godot.Resource, error) {
	if tscn.Key != TscnTypeGodotResource {
		return nil, fmt.Errorf("can't convert %s to gd_resource", tscn.Key)
	}

	res := &godot.Resource{
		ExtResources: make(map[int64]*godot.ExtResource),
		SubResources: make(map[int64]*godot.SubResource),
		Fields:       make(map[string]interface{}),
		MetaData: godot.MetaData{
			LexerPosition: tscn.Pos,
		},
	}

	t, err := tscn.GetAttribute("type")
	if err != nil {
		return nil, errors.Wrap(err, "gd_resource doesn't contain required attribute type")
	}
	if t.String == nil {
		return nil, errors.New("gd_resource attribute type must be a string")
	}

	res.Type = *t.String

	// handle everything that isn't a node
	for _, section := range tscn.Sections {
		// External resources
		if section.ResourceType == parser.ResourceTypeExtResource {
			r, err := convertSectionToExtResource(section)
			if err != nil {
				return nil, err
			}
			res.ExtResources[r.ID] = r
			continue
		}

		// Internal resources
		if section.ResourceType == parser.ResourceTypeSubResource {
			r, err := convertSectionToSubResource(section)
			if err != nil {
				return nil, err
			}
			res.SubResources[r.ID] = r
			continue
		}

		// resource section
		if section.ResourceType == parser.ResourceTypeResource {
			for _, field := range section.Fields {
				res.Fields[field.Key] = convertGdValue(field.Value)
			}
			continue
		}

		// something else found? Whoops, throw error
		return nil, fmt.Errorf("invalid resource type found: %s [%s]", section.ResourceType, section.Pos)
	}

	return res, nil
}
