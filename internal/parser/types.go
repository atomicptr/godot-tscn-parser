package parser

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

// Supported resource types
const (
	ResourceTypeExtResource = "ext_resource"
	ResourceTypeSubResource = "sub_resource"
	ResourceTypeEditable    = "editable"
	ResourceTypeConnection  = "connection"
	ResourceTypeNode        = "node"
)

// TscnFile represents a .tscn file
type TscnFile struct {
	Key        string        `parser:"('[' @Ident"`
	Attributes []*GdField    `parser:"@@* ']')?"`
	Fields     []*GdField    `parser:"@@*"`
	Sections   []*GdResource `parser:"@@*"`
	Pos        lexer.Position
}

// GetAttribute finds an attribute by name, returns error if not found
func (tscn *TscnFile) GetAttribute(name string) (*GdValue, error) {
	for _, attr := range tscn.Attributes {
		if attr.Key == name {
			return attr.Value, nil
		}
	}
	return nil, fmt.Errorf("unknown attribute in %s: %s", tscn.Pos, name)
}

// GdResource represents a resource within a TSCN file
type GdResource struct {
	ResourceType string     `parser:"'[' @Ident "`
	Attributes   []*GdField `parser:"@@* ']'"`
	Fields       []*GdField `parser:"@@*"`
	Pos          lexer.Position
}

// GetAttribute finds an attribute by name, returns error if not found
func (res *GdResource) GetAttribute(name string) (*GdValue, error) {
	for _, attr := range res.Attributes {
		if attr.Key == name {
			return attr.Value, nil
		}
	}
	return nil, fmt.Errorf("unknown attribute in %s: %s", res.Pos, name)
}

// GetField finds an attribute by name, returns error if not found
func (res *GdResource) GetField(name string) (*GdValue, error) {
	for _, field := range res.Fields {
		if field.Key == name {
			return field.Value, nil
		}
	}
	return nil, fmt.Errorf("unknown field in %s: %s", res.Pos, name)
}

// GdType represents a type with values
type GdType struct {
	Key        string     `parser:" @Ident ('('"`
	Parameters []*GdValue `parser:"(@@ ( ',' @@ )* )? ')')?"`
	Pos        lexer.Position
}

// GdField represents a field with a value
type GdField struct {
	Key   string   `parser:"@Ident '='"`
	Value *GdValue `parser:" @@"`
	Pos   lexer.Position
}

// GdMapField represents a field with a value within a map
type GdMapField struct {
	Key   string   `parser:"@String ':'"`
	Value *GdValue `parser:" @@"`
	Pos   lexer.Position
}

// ToString returns a string representation of a GdMapField
func (kv *GdMapField) ToString() string {
	return fmt.Sprintf("\"%s\": %s", kv.Key, kv.Value.ToString())
}

// GdValue represents a value
type GdValue struct {
	IsEmptyMap   *bool         `parser:" (@'{' '}')"`
	Map          []*GdMapField `parser:"| '{' ( @@ ( ',' @@ )* )? '}'"`
	KeyValuePair *GdMapField   `parser:"| @@"`
	IsEmptyArray *bool         `parser:"| (@'[' ']')"`
	Array        []*GdValue    `parser:"| '[' ( @@ ( ',' @@ )* )? (',')? ']'"`
	String       *string       `parser:"| @String"`
	Integer      *int64        `parser:"| @Int"`
	Float        *float64      `parser:"| @Float"`
	Bool         *bool         `parser:"| (@'true' | 'false')"`
	Null         *bool         `parser:"| (@'null')"`
	Type         *GdType       `parser:"| @@"`
	Pos          lexer.Position
}

// Raw returns an interface{} which contains the actual value of the associated GdValue
func (v *GdValue) Raw() interface{} {
	if (v.IsEmptyMap != nil && *v.IsEmptyMap) || len(v.Map) != 0 {
		return v.Map
	}

	if v.KeyValuePair != nil {
		return *v.KeyValuePair
	}

	if (v.IsEmptyArray != nil && *v.IsEmptyArray) || len(v.Array) != 0 {
		return v.Array
	}

	if v.String != nil {
		return *v.String
	}

	if v.Integer != nil {
		return *v.Integer
	}

	if v.Float != nil {
		return *v.Float
	}

	if v.Bool != nil {
		return *v.Bool
	}

	if v.Null != nil {
		return nil
	}

	if v.Type != nil {
		return *v.Type
	}

	return nil
}

// ToString returns a string representation of the associated GdValue
func (v *GdValue) ToString() string {
	switch value := v.Raw().(type) {
	case []*GdMapField:
		var values []string
		for _, kv := range value {
			values = append(values, kv.ToString())
		}
		return fmt.Sprintf("{%s}", strings.Join(values, ", "))
	case GdMapField:
		return value.ToString()
	case []*GdValue:
		var values []string
		for _, v := range value {
			values = append(values, v.ToString())
		}
		return fmt.Sprintf("[%s]", strings.Join(values, ", "))
	case string:
		return fmt.Sprintf("\"%s\"", value)
	case int64:
		return fmt.Sprintf("%d", value)
	case float64:
		return fmt.Sprintf("%f", value)
	case bool:
		return fmt.Sprintf("%v", value)
	case GdType:
		var values []string
		for _, param := range value.Parameters {
			values = append(values, param.ToString())
		}
		return fmt.Sprintf("%s (%s)", v.Type.Key, strings.Join(values, ", "))
	default:
		return "null"
	}
}
