// Package parser implements a participle parser for the TSCN file format
package parser

import (
	"io"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer/stateful"
)

var tscnLexer = stateful.MustSimple([]stateful.Rule{
	{
		Name:    "Ident",
		Pattern: `([0-9]+)?/?[a-zA-Z_][a-zA-Z_\d/]*`,
		Action:  nil,
	},
	{
		Name:    "String",
		Pattern: `"[^"]*"`,
		Action:  nil,
	},
	{
		Name:    "Float",
		Pattern: `-?[0-9]*\.[0-9]+(e\-[0-9]+)?\b`,
		Action:  nil,
	},
	{
		Name:    "Int",
		Pattern: `-?[0-9]+\b`,
		Action:  nil,
	},
	{
		Name:    "Punct",
		Pattern: `[][=():,{}]`,
		Action:  nil,
	},
	{
		Name:    "comment",
		Pattern: `;[^\n]*`,
		Action:  nil,
	},
	{
		Name:    "whitespace",
		Pattern: `\s+`,
		Action:  nil,
	},
})

var tscnParser = participle.MustBuild(
	&TscnFile{},
	participle.Lexer(tscnLexer),
	participle.Unquote("String"),
)

// Parse content and return a simple representation of the file format.
func Parse(r io.Reader) (*TscnFile, error) {
	ast := &TscnFile{}
	err := tscnParser.Parse("", r, ast)
	if err != nil {
		return nil, err
	}
	return ast, nil
}
