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

// TODO: possibly not a very useful function
func (v GdValue) ToString() string {
	if v.String != nil {
		return *v.String
	}

	if v.Integer != nil {
		return fmt.Sprintf("%d", *v.Integer)
	}

	if v.Float != nil {
		return fmt.Sprintf("%f", *v.Float)
	}

	if v.Bool != nil {
		return fmt.Sprintf("%v", *v.Bool)
	}

	if len(v.Map) != 0 {
		var values []string
		for _, kv := range v.Map {
			values = append(values, fmt.Sprintf("\"%s\": %s", kv.Key, kv.Value.ToString()))
		}
		return fmt.Sprintf("Map {%s}", strings.Join(values, ", "))
	}

	if len(v.Array) != 0 {
		var values []string
		for _, value := range v.Array {
			values = append(values, value.ToString())
		}
		return fmt.Sprintf("[%s]", strings.Join(values, ", "))
	}

	if v.Null != nil {
		return "null"
	}

	if v.Type != nil {
		var values []string

		for _, param := range v.Type.Parameters {
			values = append(values, param.ToString())
		}

		return fmt.Sprintf("%s (%s)", v.Type.Key, strings.Join(values, ", "))
	}

	return "???"
}
