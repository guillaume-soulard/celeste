package model

import (
	"celeste/src/model/ast"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

func Parse(str string) (parsed *ast.Ast, err error) {
	basicLexer := lexer.MustSimple([]lexer.SimpleRule{
		{"String", `'(\\'|[^'])*'`},
		{"JsonString", `"(\\"|[^"])*"`},
		{"Number", `[-+]?(\d*\.)?\d+`},
		{"Ident", `[a-zA-Z_]\w*`},
		{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`},
		{"whitespace", `\s+`},
	})
	basicParser := participle.MustBuild[ast.Ast](
		participle.Lexer(basicLexer),
		participle.CaseInsensitive("Ident"),
		participle.Unquote("String"),
		participle.UseLookahead(2),
	)
	parsed, err = basicParser.ParseString("", str)
	return parsed, err
}
