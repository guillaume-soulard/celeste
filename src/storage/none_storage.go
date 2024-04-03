package storage

import (
	"celeste/src/model"
	"celeste/src/model/ast"
	"time"
)

func NewNoneStorage() Storage {
	return &None{IdGenerator: NewIdGenerator()}
}

type None struct {
	IdGenerator IdGenerator
}

func (n *None) Append(_ ast.Json) (id string, err error) {
	id = n.IdGenerator.NextId(time.Now())
	return id, err
}

func (n *None) InitCursor(_ model.StartPosition) (cursor interface{}, err error) {
	cursor = 0
	return cursor, err
}

func (n *None) Read(_ model.ReadBehaviour, cursor interface{}, _ int) (newCursor interface{}, data []ast.Json, endOfStream bool, err error) {
	data = make([]ast.Json, 0)
	endOfStream = true
	newCursor = cursor
	return newCursor, data, endOfStream, err
}

func (n *None) Close() (err error) {
	return err
}

func (n *None) Truncate(_ *[]ast.EvictionPolicy) (err error) {
	return err
}
