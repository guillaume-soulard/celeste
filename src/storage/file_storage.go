package storage

import (
	"celeste/src/model"
	"celeste/src/model/ast"
	"fmt"
	"io"
	"os"
	"time"
	"unsafe"
)

func NewFileStorage(streamName string) (storage Storage, err error) {
	var files []*os.File
	if files, err = openDatabaseFiles(streamName); err != nil {
		return storage, err
	}
	storage = &File{
		files:                 files,
		streamName:            streamName,
		PreviousReadBehaviour: model.ReadBehaviourNext,
	}
	return storage, err
}

func openDatabaseFiles(streamName string) (files []*os.File, err error) {
	var streamDir *os.File
	if _, err = os.Stat(streamName); os.IsNotExist(err) {
		if err = os.Mkdir(streamName, 0644); err != nil {
			return files, err
		}
	}
	if streamDir, err = os.OpenFile(streamName, os.O_CREATE, 0644); err != nil {
		return files, err
	}
	var allFiles []os.FileInfo
	if allFiles, err = streamDir.Readdir(-1); err != nil {
		return files, err
	}
	files = make([]*os.File, 0, 10)
	for _, f := range allFiles {
		if !f.IsDir() {
			var file *os.File
			if file, err = os.OpenFile(fmt.Sprintf("%s/%s", streamName, f.Name()), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644); err != nil {
				return files, err
			}
			files = append(files, file)
		}
	}
	return files, err
}

type File struct {
	files                 []*os.File
	size                  uint64
	len                   uint64
	streamName            string
	oldestInsertedTime    time.Time
	PreviousReadBehaviour model.ReadBehaviour
}

const sizeLen = 4

func (f *File) Append(data ast.Json) (id int64, err error) {
	var file *os.File
	if file, err = f.GetLastFile(); err != nil {
		return id, err
	}
	if id, err = file.Seek(0, io.SeekCurrent); err != nil {
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
	_, err = file.Write(b)
	return id, err
}

func (f *File) GetLastFile() (file *os.File, err error) {
	if len(f.files) == 0 {
		file, err = os.OpenFile(fmt.Sprintf("%s/%d.dat", f.streamName, len(f.files)), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	} else {
		file = f.files[len(f.files)-1]
	}
	return file, err
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
	if f.files != nil {
		for _, file := range f.files {
			if err = file.Close(); err != nil {
				return err
			}
			file = nil
		}
	}
	return err
}

type FileCursor struct {
	inFileCursor int64
	fileIndex    int
}

func (f *File) InitCursor(startPosition model.StartPosition) (cursor interface{}, err error) {
	var file *os.File
	if file, err = f.GetLastFile(); err != nil {
		return cursor, err
	}
	var inFileCursor int64
	var fileIndex int
	if startPosition == model.StartPositionBeginning {
		inFileCursor = int64(0)
		fileIndex = 0
	} else if startPosition == model.StartPositionEnd {
		if inFileCursor, err = file.Seek(0, io.SeekEnd); err != nil {
			return cursor, err
		}
		fileIndex = len(f.files) - 1
	}
	cursor = FileCursor{
		inFileCursor: inFileCursor,
		fileIndex:    fileIndex,
	}
	return cursor, err
}

func (f *File) Read(readBehaviour model.ReadBehaviour, cursor interface{}, count int) (newCursor interface{}, data []ast.Json, endOfStream bool, err error) {
	fc := cursor.(FileCursor)
	file := f.files[fc.fileIndex]
	if _, err = file.Seek(cursor.(int64), io.SeekStart); err != nil {
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

func (f *File) readBackWard(data *[]ast.Json, newCursor *interface{}) (endOfFile bool, err error) {
	fc := (*newCursor).(FileCursor)
	file := f.files[fc.fileIndex]
	sizeBytes := make([]byte, 4)
	if _, endOfFile, err = f.TrySeek(file, -int64(sizeLen), io.SeekCurrent); err != nil || endOfFile {
		return endOfFile, err
	}
	if endOfFile, err = readSize(file, &sizeBytes); err != nil || endOfFile {
		return endOfFile, err
	}
	size := ByteArrayToInt(sizeBytes)
	if _, endOfFile, err = f.TrySeek(file, -(size + sizeLen), io.SeekCurrent); err != nil || endOfFile {
		return endOfFile, err
	}
	if _, endOfFile, err = readData(file, size, data); err != nil {
		return endOfFile, err
	}
	var cursor int64
	if cursor, endOfFile, err = f.TrySeek(file, -(size + sizeLen), io.SeekCurrent); err != nil || endOfFile {
		return endOfFile, err
	}
	fileIndex := fc.fileIndex
	if endOfFile {
		fileIndex++
	}
	*newCursor = FileCursor{
		inFileCursor: cursor,
		fileIndex:    fileIndex,
	}
	return endOfFile, err
}

func (f *File) TrySeek(file *os.File, offset int64, whence int) (cursor int64, endOfFile bool, err error) {
	var currentCursorPosition int64
	if currentCursorPosition, err = file.Seek(0, io.SeekCurrent); err != nil {
		return cursor, endOfFile, err
	}
	if currentCursorPosition+offset < 0 {
		endOfFile = true
		return cursor, endOfFile, err
	}
	if cursor, err = file.Seek(offset, whence); err != nil {
		if err == io.EOF {
			err = nil
			endOfFile = true
		}
	}
	return cursor, endOfFile, err
}

func (f *File) readForward(data *[]ast.Json, newCursor *interface{}) (endOfFile bool, err error) {
	fc := (*newCursor).(FileCursor)
	file := f.files[fc.fileIndex]
	sizeBytes := make([]byte, 4)
	var nbRead int
	if endOfFile, err = readSize(file, &sizeBytes); err != nil || endOfFile {
		return endOfFile, err
	}
	size := ByteArrayToInt(sizeBytes)
	if nbRead, endOfFile, err = readData(file, size, data); err != nil || endOfFile {
		return endOfFile, err
	}
	if endOfFile, err = readSize(file, &sizeBytes); err != nil || endOfFile {
		return endOfFile, err
	}
	fileIndex := fc.fileIndex
	if endOfFile {
		fileIndex++
	}
	*newCursor = FileCursor{
		inFileCursor: fc.inFileCursor + int64(4+nbRead),
		fileIndex:    fileIndex,
	}
	return endOfFile, err
}

func readData(f *os.File, size int64, data *[]ast.Json) (nbRead int, endOfStream bool, err error) {
	byteData := make([]byte, size)
	if nbRead, err = f.Read(byteData); err != nil {
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

func readSize(f *os.File, size *[]byte) (endOfStream bool, err error) {
	if _, err = f.Read(*size); err != nil {
		if err == io.EOF {
			err = nil
			endOfStream = true
		}
	}
	return endOfStream, err
}

func (f *File) Truncate(evictionPolicies *[]ast.EvictionPolicy) (err error) {
	for _, policy := range *evictionPolicies {
		if policy.MaxAmountItems != nil && uint64(*policy.MaxAmountItems) < f.len {
			if err = f.truncateByMaxItems(policy); err != nil {
				return err
			}
		}
		if policy.MaxSize != nil && (*policy.MaxSize).Bytes() < f.size {
			if err = f.truncateByMaxSize(policy); err != nil {
				return err
			}
		}
		if policy.MaxDuration != nil && time.Now().Sub(f.oldestInsertedTime) > (*policy.MaxDuration).Duration() {
			if err = f.truncateByMaxDuration(policy); err != nil {
				return err
			}
		}
	}
	return err
}

func (f *File) truncateByMaxItems(policy ast.EvictionPolicy) (err error) {
	return err
}

func (f *File) truncateByMaxSize(policy ast.EvictionPolicy) (err error) {
	return err
}

func (f *File) truncateByMaxDuration(policy ast.EvictionPolicy) (err error) {
	return err
}
