package database

import (
	"celeste/src/model/ast"
)

type StreamDropExecutor struct{}

func (s StreamDropExecutor) Accept(ast *ast.Ast) bool {
	return ast.StreamDrop != nil
}

func (s StreamDropExecutor) Execute(db *Database, ast *ast.Ast) (err error) {
	err = db.DropStream(*ast.StreamDrop.Name)
	return err
}
