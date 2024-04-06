package database

import (
	"celeste/src/model/ast"
)

type StreamAppendExecutor struct{}

func (s StreamAppendExecutor) Accept(ast *ast.Ast) bool {
	return ast.StreamAppend != nil
}

func (s StreamAppendExecutor) Execute(db *Database, ast *ast.Ast) (err error) {
	streams := make([]*Stream, len(*ast.StreamAppend.StreamNames))
	for i, streamName := range *ast.StreamAppend.StreamNames {
		if streams[i], err = db.LookUpStream(streamName); err != nil {
			return err
		}
	}
	ids := make([]string, len(*ast.StreamAppend.JsonData)*len(*ast.StreamAppend.StreamNames))
	var id string
	index := 0
	for _, stream := range streams {
		for _, data := range *ast.StreamAppend.JsonData {
			if id, err = stream.Append(data); err != nil {
				return err
			}
			ids[index] = id
			index++
		}
	}
	return err
}
