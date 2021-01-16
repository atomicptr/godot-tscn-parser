package parser

import (
	"fmt"
	"github.com/alecthomas/participle/v2/lexer"
	"strings"
)

type GdField struct {
	Key   string   `parser:"@Ident \"=\""`
	Value *GdValue `parser:" @@"`
	Pos   lexer.Position
}

type GdMapField struct {
	Key   string   `parser:"@String \":\""`
	Value *GdValue `parser:" @@"`
	Pos   lexer.Position
}

func (kv GdMapField) ToString() string {
	return fmt.Sprintf("\"%s\": %s", kv.Key, kv.Value.ToString())
}

type GdValue struct {
	Map          []*GdMapField `parser:" \"{\" ( @@ ( \",\" @@ )* )? \"}\""`
	KeyValuePair *GdMapField   `parser:"| @@"`
	Array        []*GdValue    `parser:"| \"[\" ( @@ ( \",\" @@ )* )? (\",\")? \"]\""`
	String       *string       `parser:"| @String"`
	Integer      *int64        `parser:"| @Int"`
	Float        *float64      `parser:"| @Float"`
	Bool         *bool         `parser:"| (@\"true\" | \"false\")"`
	Null         *bool         `parser:"| (@\"null\")"`
	Type         *GdType       `parser:"| @@"`
	Pos          lexer.Position
}

func (v GdValue) Raw() interface{} {
	if len(v.Map) != 0 {
		return v.Map
	}

	if v.KeyValuePair != nil {
		return *v.KeyValuePair
	}

	if len(v.Array) != 0 {
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

func (v GdValue) ToString() string {
	switch value := v.Raw().(type) {
	case []*GdMapField:
		var values []string
		for _, kv := range value {
			values = append(values, kv.ToString())
		}
		return fmt.Sprintf("Map {%s}", strings.Join(values, ", "))
	case GdMapField:
		return value.ToString()
	case []*GdValue:
		var values []string
		for _, v := range value {
			values = append(values, v.ToString())
		}
		return fmt.Sprintf("[%s]", strings.Join(values, ", "))
	case string:
		return value
	case int64:
		return fmt.Sprintf("%d", value)
	case float64:
		return fmt.Sprintf("%f", value)
	case bool:
		return fmt.Sprintf("%v", value)
	case nil:
		return "null"
	case GdType:
		var values []string
		for _, param := range value.Parameters {
			values = append(values, param.ToString())
		}
		return fmt.Sprintf("%s (%s)", v.Type.Key, strings.Join(values, ", "))
	}

	return "???"
}
