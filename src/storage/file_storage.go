package storage

import (
	"celeste/src/model"
	"celeste/src/model/ast"
	"errors"
	"io"
	"os"
	"unsafe"
)

func NewFileStorage(streamName string) (storage Storage, err error) {
	var file *os.File
	file, err = os.OpenFile(streamName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	storage = &File{
		file:                  file,
		PreviousReadBehaviour: model.ReadBehaviourNext,
	}
	return storage, err
}

type File struct {
	file                  *os.File
	PreviousReadBehaviour model.ReadBehaviour
}

const sizeLen = 4

func (f *File) Append(data ast.Json) (id int64, err error) {
	if id, err = f.file.Seek(0, io.SeekCurrent); err != nil {
		return id, err
	}
	dataStr := data.ToString()
	strLen := len(dataStr)
	b := make([]byte, (sizeLen*2)+strLen)
	sizeByte := IntToByteArray(strLen)
	str := []byte(dataStr)
	copy(b[0:sizeLen], sizeByte)
	copy(b[sizeLen:len(str)+sizeLen], str)
	copy(b[sizeLen+len(str):], sizeByte)
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
	data = make([]ast.Json, 0, count)
	newCursor = cursor.(int64)
	for !endOfStream && len(data) < count {
		if readBehaviour == model.ReadBehaviourNext || (readBehaviour == model.ReadBehaviourAgain && f.PreviousReadBehaviour == model.ReadBehaviourNext) {
			if endOfStream, err = f.readForward(&data, &newCursor); err != nil || endOfStream {
				return newCursor, data, endOfStream, err
			}
		} else if readBehaviour == model.ReadBehaviourPrevious || (readBehaviour == model.ReadBehaviourAgain && f.PreviousReadBehaviour == model.ReadBehaviourPrevious) {
			if endOfStream, err = f.readBackWard(&data, &newCursor); err != nil {
				return newCursor, data, endOfStream, err
			}
		}
	}
	f.PreviousReadBehaviour = readBehaviour
	return newCursor, data, endOfStream, err
}

func (f *File) readBackWard(data *[]ast.Json, newCursor *interface{}) (endOfStream bool, err error) {
	sizeBytes := make([]byte, 4)
	if _, endOfStream, err = f.TrySeek(-int64(sizeLen), io.SeekCurrent); err != nil || endOfStream {
		return endOfStream, err
	}
	if endOfStream, err = readSize(f, &sizeBytes); err != nil || endOfStream {
		return endOfStream, err
	}
	size := ByteArrayToInt(sizeBytes)
	if _, endOfStream, err = f.TrySeek(-(size + sizeLen), io.SeekCurrent); err != nil || endOfStream {
		return endOfStream, err
	}
	if _, endOfStream, err = readData(f, size, data); err != nil {
		return endOfStream, err
	}
	var cursor int64
	if cursor, endOfStream, err = f.TrySeek(-(size + sizeLen), io.SeekCurrent); err != nil || endOfStream {
		return endOfStream, err
	}
	*newCursor = cursor
	return endOfStream, err
}

func (f *File) TrySeek(offset int64, whence int) (cursor int64, endOfStream bool, err error) {
	var currentCursorPosition int64
	if currentCursorPosition, err = f.file.Seek(0, io.SeekCurrent); err != nil {
		return cursor, endOfStream, err
	}
	if currentCursorPosition+offset < 0 {
		endOfStream = true
		return cursor, endOfStream, err
	}
	if cursor, err = f.file.Seek(offset, whence); err != nil {
		if err == io.EOF {
			err = nil
			endOfStream = true
		}
	}
	return cursor, endOfStream, err
}

func (f *File) readForward(data *[]ast.Json, newCursor *interface{}) (endOfStream bool, err error) {
	sizeBytes := make([]byte, 4)
	var nbRead int
	if endOfStream, err = readSize(f, &sizeBytes); err != nil || endOfStream {
		return endOfStream, err
	}
	size := ByteArrayToInt(sizeBytes)
	if nbRead, endOfStream, err = readData(f, size, data); err != nil || endOfStream {
		return endOfStream, err
	}
	if endOfStream, err = readSize(f, &sizeBytes); err != nil || endOfStream {
		return endOfStream, err
	}
	*newCursor = (*newCursor).(int64) + int64(4+nbRead)
	return endOfStream, err
}

func readData(f *File, size int64, data *[]ast.Json) (nbRead int, endOfStream bool, err error) {
	byteData := make([]byte, size)
	if nbRead, err = f.file.Read(byteData); err != nil {
		if err == io.EOF {
			err = nil
			endOfStream = true
		}
		return nbRead, endOfStream, err
	}
	var parsed *ast.Ast
	if parsed, err = model.Parse(string(byteData)); err != nil {
		return nbRead, endOfStream, err
	}
	if parsed.Json != nil {
		*data = append(*data, *parsed.Json)
	}
	return nbRead, endOfStream, err
}

func readSize(f *File, size *[]byte) (endOfStream bool, err error) {
	if _, err = f.file.Read(*size); err != nil {
		if err == io.EOF {
			err = nil
			endOfStream = true
		}
	}
	return endOfStream, err
}

func (f *File) Truncate(_ *[]ast.EvictionPolicy) (err error) {
	err = errors.New("not implemented")
	return err
}
