package parser

import "github.com/atomicptr/godot-tscn-parser/pkg/godot"

const (
	// TscnTypeGodotScene is the identifier for a Godot Scene
	TscnTypeGodotScene = "gd_scene"
)

func insertFieldEntriesFromSection(section *GdResource, fieldMap map[string]interface{}) {
	for _, field := range section.Fields {
		// TODO: properly parse structures like bones/0/name = "Bone", bones/0/parent = -1, etc.
		fieldMap[field.Key] = convertGdValue(field.Value)
	}
}

func convertGdValue(val *GdValue) interface{} {
	switch value := val.Raw().(type) {
	case []*GdValue:
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
	case []*GdMapField:
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
	case GdType:
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
	case GdMapField:
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
