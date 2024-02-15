package storage

import (
	"celeste/src/model"
	"celeste/src/model/ast"
	"io"
	"os"
)

func NewFileStorage(streamName string) (storage Storage, err error) {
	var file *os.File
	file, err = os.OpenFile(streamName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	storage = &File{
		file: file,
	}
	return storage, err
}

type File struct {
	file *os.File
}

func (f *File) Append(data ast.Json) (id int64, err error) {
	if id, err = f.file.Seek(0, io.SeekCurrent); err != nil {
		return id, err
	}
	_, err = f.file.WriteString("test")
	return id, err
}

func (f *File) InitCursor(startPosition model.StartPosition) (cursor interface{}, err error) {
	panic("implement me")
}

func (f *File) Read(readBehaviour model.ReadBehaviour, cursor interface{}, count int) (newCursor interface{}, data []ast.Json, endOfStream bool, err error) {
	panic("implement me")
}
