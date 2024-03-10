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
	Close() (err error)
	Truncate(evictionPolicies *[]ast.EvictionPolicy) (err error)
}

func NewStorageFrom(streamCreation ast.StreamCreation) (s Storage, err error) {
	if streamCreation.Storage == nil || streamCreation.Storage.Disk {
		s, err = NewFileStorage(*streamCreation.Name)
	} else if streamCreation.Storage.Memory {
		s = NewMemoryStorage()
	} else if streamCreation.Storage.None {
		s = NewNoneStorage()
	} else {
		err = errors.New("unknown storage")
	}
	return s, err
}
