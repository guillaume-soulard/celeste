package storage

import (
	"celeste/src/model"
	"celeste/src/model/ast"
	"errors"
)

type Storage interface {
	Append(data ast.Json) (id int64, err error)
	InitCursor(position model.StartPosition) (cursor interface{}, err error)
	Read(db model.ReadBehaviour, cursor interface{}, count int) (newCursor interface{}, data []ast.Json, endOfStream bool, err error)
}

func NewStorageFrom(storage *ast.StreamStorage) (s Storage, err error) {
	if storage.Disk {
		s = NewFileStorage()
	} else if storage.Memory {
		s = NewMemoryStorage()
	} else if storage.None {
		s = NewNoneStorage()
	} else {
		err = errors.New("unknown storage")
	}
	return s, err
}
