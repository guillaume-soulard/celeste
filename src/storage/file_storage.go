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
	file, err = os.OpenFile(streamName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
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

func ByteArrayToInt(arr []byte) int64 {
	val := int64(0)
	size := len(arr)
	for i := 0; i < size; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
	}
	return val
}

func (f *File) Close() (err error) {
	if f.file != nil {
		err = f.file.Close()
		f.file = nil
	}
	return err
}

func (f *File) InitCursor(startPosition model.StartPosition) (cursor interface{}, err error) {
	if startPosition == model.StartPositionBeginning {
		cursor = int64(0)
	} else if startPosition == model.StartPositionEnd {
		cursor, err = f.file.Seek(0, io.SeekEnd)
	}
	return cursor, err
}

func (f *File) Read(readBehaviour model.ReadBehaviour, cursor interface{}, count int) (newCursor interface{}, data []ast.Json, endOfStream bool, err error) {
	if _, err = f.file.Seek(cursor.(int64), io.SeekStart); err != nil {
		return newCursor, data, endOfStream, err
	}
	size := make([]byte, 4)
	data = make([]ast.Json, 0, count)
	var nbRead int
	newCursor = cursor.(int64)
	for len(data) < count {
		if _, err = f.file.Read(size); err != nil {
			if err == io.EOF {
				err = nil
				endOfStream = true
			}
			return newCursor, data, endOfStream, err
		}
		byteData := make([]byte, ByteArrayToInt(size))
		if nbRead, err = f.file.Read(byteData); err != nil {
			if err == io.EOF {
				err = nil
				endOfStream = true
			}
			return newCursor, data, endOfStream, err
		}
		var parsed *ast.Ast
		if parsed, err = model.Parse(string(byteData)); err != nil {
			return newCursor, data, endOfStream, err
		}
		if parsed.Json != nil {
			data = append(data, *parsed.Json)
		}
		newCursor = newCursor.(int64) + int64(4+nbRead)
	}
	return newCursor, data, endOfStream, err
}
