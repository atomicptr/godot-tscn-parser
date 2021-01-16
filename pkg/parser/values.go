package parser

import (
	"fmt"
	"github.com/alecthomas/participle/v2/lexer"
	"strings"
)

type GdField struct {
	Key   string   `@Ident "="`
	Value *GdValue ` @@`
	Pos   lexer.Position
}

type GdMapField struct {
	Key   string   `@String ":"`
	Value *GdValue ` @@`
	Pos   lexer.Position
}

type GdValue struct {
	String  *string       ` @String`
	Integer *int64        `| @Int`
	Float   *float64      `| @Float`
	Bool    *bool         `| (@"true" | "false")`
	Map     []*GdMapField `| "{" ( @@ ( "," @@ )* )? "}"`
	Array   []*GdValue    `| "[" ( @@ ( "," @@ )* )? (",")? "]"`
	Null    *bool         `| (@"null")`
	Type    *GdType       `| @@`
	Pos     lexer.Position
}

func (v GdValue) Raw() interface{} {
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

	if len(v.Map) != 0 {
		return v.Map
	}

	if len(v.Array) != 0 {
		return v.Array
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
	case string:
		return value
	case int64:
		return fmt.Sprintf("%d", value)
	case float64:
		return fmt.Sprintf("%f", value)
	case bool:
		return fmt.Sprintf("%v", value)
	case []*GdMapField:
		var values []string
		for _, kv := range value {
			values = append(values, fmt.Sprintf("\"%s\": %s", kv.Key, kv.Value.ToString()))
		}
		return fmt.Sprintf("Map {%s}", strings.Join(values, ", "))
	case []*GdValue:
		var values []string
		for _, v := range value {
			values = append(values, v.ToString())
		}
		return fmt.Sprintf("[%s]", strings.Join(values, ", "))
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
