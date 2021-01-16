package parser

import "github.com/alecthomas/participle/v2/lexer"

type TscnFile struct {
	Key        string        `parser:"(\"[\" @Ident"`
	Attributes []*GdField    `parser:"@@* \"]\")?"`
	Fields     []*GdField    `parser:"@@*"`
	Sections   []*GdResource `parser:"@@*"`
	Pos        lexer.Position
}

type GdResource struct {
	ResourceType string     `parser:"\"[\" @Ident "`
	Attributes   []*GdField `parser:"@@* \"]\""`
	Fields       []*GdField `parser:"@@*"`
	Pos          lexer.Position
}

type GdType struct {
	Key        string     `parser:" @Ident (\"(\""`
	Parameters []*GdValue `parser:"@@ ( \",\" @@ )* \")\")?"`
	Pos        lexer.Position
}
