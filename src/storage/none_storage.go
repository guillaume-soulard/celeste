package storage

import (
	"celeste/src/model"
	"celeste/src/model/ast"
)

func NewNoneStorage() Storage {
	return &None{LastId: 0}
}

type None struct {
	LastId int64
}

func (n *None) Append(_ ast.Json) (id int64, err error) {
	n.LastId++
	id = n.LastId
	return id, err
}

func (n *None) InitCursor(_ model.StartPosition) (cursor interface{}, err error) {
	cursor = 0
	return cursor, err
}

func (n *None) Read(_ model.ReadBehaviour, _ interface{}, _ int) (newCursor interface{}, data []ast.Json, endOfStream bool, err error) {
	data = make([]ast.Json, 0)
	endOfStream = true
	return newCursor, data, endOfStream, err
}
