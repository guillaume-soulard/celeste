package database

import (
	"celeste/src/model/ast"
)

type StreamReadExecutor struct{}

func (s StreamReadExecutor) Accept(ast *ast.Ast) bool {
	return ast.StreamRead != nil
}

func (s StreamReadExecutor) Execute(db *Database, ast *ast.Ast) (err error) {
	read := ast.StreamRead
	var stream *Stream
	if stream, err = db.LookUpStream(*read.ReadClause.Stream); err != nil {
		return err
	}
	var iterator StreamIterator
	if iterator, err = stream.Iterator(ast.StreamRead); err != nil {
		return err
	}
	err = iterator.Read(ast.StreamRead)
	return err
}
