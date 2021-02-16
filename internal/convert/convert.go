// Package convert deals with converting parser output into an usable data structure see pkg/godot
package convert

import (
	"github.com/atomicptr/godot-tscn-parser/internal/parser"
	"github.com/atomicptr/godot-tscn-parser/pkg/godot"
)

const (
	// TscnTypeGodotScene is the identifier for a Godot Scene
	TscnTypeGodotScene = "gd_scene"
	// TscnTypeGodotResource is the identifier for .tres files
	TscnTypeGodotResource = "gd_resource"
)

func insertFieldEntriesFromSection(section *parser.GdResource, fieldMap map[string]interface{}) {
	for _, field := range section.Fields {
		// TODO: properly parse structures like bones/0/name = "Bone", bones/0/parent = -1, etc.
		fieldMap[field.Key] = convertGdValue(field.Value)
	}
}

func convertGdValue(val *parser.GdValue) interface{} {
	switch value := val.Raw().(type) {
	case []*parser.GdValue:
		values := make([]interface{}, len(value))
		for index, v := range value {
			values[index] = v
		}
		return godot.Value{
			Value: values,
			MetaData: godot.MetaData{
				LexerPosition: val.Pos,
			},
		}
	case []*parser.GdMapField:
		m := make(map[string]interface{})
		for _, kv := range value {
			m[kv.Key] = convertGdValue(kv.Value)
		}
		return godot.Value{
			Value: m,
			MetaData: godot.MetaData{
				LexerPosition: val.Pos,
			},
		}
	case parser.GdType:
		params := make([]interface{}, len(value.Parameters))
		for index, p := range value.Parameters {
			params[index] = convertGdValue(p)
		}
		return godot.Type{
			Identifier: value.Key,
			Parameters: params,
			MetaData: godot.MetaData{
				LexerPosition: value.Pos,
			},
		}
	case parser.GdMapField:
		return godot.KeyValuePair{
			Key:   value.Key,
			Value: convertGdValue(value.Value),
			MetaData: godot.MetaData{
				LexerPosition: value.Pos,
			},
		}
	default:
		return godot.Value{
			Value: value,
			MetaData: godot.MetaData{
				LexerPosition: val.Pos,
			},
		}
	}
}
