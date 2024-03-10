package database

import (
	"celeste/src/model/ast"
)

type StreamTruncateExecutor struct{}

func (s StreamTruncateExecutor) Accept(ast *ast.Ast) bool {
	return ast.StreamTruncate != nil
}

func (s StreamTruncateExecutor) Execute(db *Database, ast *ast.Ast) (err error) {
	var stream *Stream
	if stream, err = db.LookUpStream(*ast.StreamTruncate.Name); err != nil {
		return err
	}
	err = stream.Storage.Truncate(ast.StreamTruncate.EvictionPolicies)
	return err
}
