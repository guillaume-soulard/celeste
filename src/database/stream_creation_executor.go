package database

import (
	"celeste/src/model/ast"
)

type StreamCreationExecutor struct{}

func (s StreamCreationExecutor) Accept(ast *ast.Ast) bool {
	return ast.StreamCreation != nil
}

func (s StreamCreationExecutor) Execute(db *Database, ast *ast.Ast) (err error) {
	err = NewStreamFrom(db, *ast.StreamCreation)
	return err
}
