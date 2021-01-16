package parser

import "github.com/alecthomas/participle/v2/lexer"

// TscnFile represents a .tscn file
type TscnFile struct {
	Key        string        `parser:"(\"[\" @Ident"`
	Attributes []*GdField    `parser:"@@* \"]\")?"`
	Fields     []*GdField    `parser:"@@*"`
	Sections   []*GdResource `parser:"@@*"`
	Pos        lexer.Position
}

// GdResource represents a resource within a TSCN file
type GdResource struct {
	ResourceType string     `parser:"\"[\" @Ident "`
	Attributes   []*GdField `parser:"@@* \"]\""`
	Fields       []*GdField `parser:"@@*"`
	Pos          lexer.Position
}

// GdType represents a type with values
type GdType struct {
	Key        string     `parser:" @Ident (\"(\""`
	Parameters []*GdValue `parser:"@@ ( \",\" @@ )* \")\")?"`
	Pos        lexer.Position
}
