package database

import (
	"celeste/src/model/ast"
)

var CommandExecutors = []CommandExecutor{
	StreamCreationExecutor{},
	StreamAppendExecutor{},
	StreamReadExecutor{},
}

type CommandExecutor interface {
	Accept(ast *ast.Ast) bool
	Execute(db *Database, ast *ast.Ast) (err error)
}
