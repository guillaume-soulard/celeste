package storage

import (
	"celeste/src/model"
	"celeste/src/model/ast"
)

func NewFileStorage() Storage {
	return &File{}
}

type File struct {
}

func (f *File) Append(data ast.Json) (id int64, err error) {
	panic("implement me")
}

func (f *File) InitCursor(startPosition model.StartPosition) (cursor interface{}, err error) {
	panic("implement me")
}

func (f *File) Read(readBehaviour model.ReadBehaviour, cursor interface{}, count int) (newCursor interface{}, data []ast.Json, endOfStream bool, err error) {
	panic("implement me")
}
