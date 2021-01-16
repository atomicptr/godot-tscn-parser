package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer/stateful"
	"io"
)

var tscnLexer = stateful.MustSimple([]stateful.Rule{
	{"Ident", `([0-9]+/)?[a-zA-Z_][a-zA-Z_\d/]*`, nil},
	{"String", `"[^"]*"`, nil},
	{"Float", `-?[0-9]*\.[0-9]+(e\-[0-9]+)?\b`, nil},
	{"Int", `-?[0-9]+\b`, nil},
	{"Punct", `[][=():,{}]`, nil},
	{"comment", `;[^\n]*`, nil},
	{"whitespace", `\s+`, nil},
})

var tscnParser = participle.MustBuild(
	&TscnFile{},
	participle.Lexer(tscnLexer),
	participle.Unquote("String"),
)

func Parse(r io.Reader) (*TscnFile, error) {
	ast := &TscnFile{}
	err := tscnParser.Parse("", r, ast)
	if err != nil {
		return nil, err
	}
	return ast, nil
}
