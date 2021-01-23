package godot

import "github.com/alecthomas/participle/v2/lexer"

// MetaData contains extra informations obtained from parsing
type MetaData struct {
	// LexerPosition the position where the entry was found in the source file
	LexerPosition lexer.Position
}
