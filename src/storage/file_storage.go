package storage

import (
	"celeste/src/model"
	"celeste/src/model/ast"
	"io"
	"os"
	"unsafe"
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
	dataStr := data.ToString()
	strLen := len(dataStr)
	b := make([]byte, 4+strLen)
	sizeByte := IntToByteArray(strLen)
	str := []byte(dataStr)
	copy(b[0:4], sizeByte)
	copy(b[4:], str)
	_, err = f.file.Write(b)
	return id, err
}

func IntToByteArray(num int) []byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	for i := 0; i < size; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr
}

func (f *File) Close() (err error) {
	if f.file != nil {
		err = f.file.Close()
	}
	return err
}

func (f *File) InitCursor(startPosition model.StartPosition) (cursor interface{}, err error) {
	panic("implement me")
}

func (f *File) Read(readBehaviour model.ReadBehaviour, cursor interface{}, count int) (newCursor interface{}, data []ast.Json, endOfStream bool, err error) {
	panic("implement me")
}
