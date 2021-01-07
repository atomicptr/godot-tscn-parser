package parser

import "github.com/alecthomas/participle/v2/lexer"

type GdScene struct {
	FileDescriptors []*GdField    `"[" "gd_scene" @@* "]"`
	Sections        []*GdResource `@@*`
	Pos             lexer.Position
}

type GdResource struct {
	ResourceType string     `"[" @Ident `
	Attributes   []*GdField `@@* "]"`
	Fields       []*GdField `@@*`
	Pos          lexer.Position
}

type GdType struct {
	Key        string     ` @Ident "("`
	Parameters []*GdValue `@@ ( "," @@ )* ")"`
	Pos        lexer.Position
}
