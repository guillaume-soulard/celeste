package storage

import (
	"celeste/src/model"
	"celeste/src/model/ast"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_FileStorage_NewFileStorage_should_return_storage(t *testing.T) {
	// GIVEN
	var storage Storage
	var err error
	streamName := "logs"
	defer func() {
		err = storage.Close()
		cleanFile(t, streamName)
	}()
	// WHEN
	storage, err = NewFileStorage(streamName)
	// THEN
	assert.NoError(t, err)
	assert.NotNil(t, storage)
}

func Test_FileStorage_Append_should_return_id_0_on_first_item(t *testing.T) {
	// GIVEN
	streamName := "logs"
	storage, err := NewFileStorage(streamName)
	defer func() {
		err = storage.Close()
		cleanFile(t, streamName)
	}()
	var id int64
	// WHEN
	id, err = storage.Append(stringToJson(`{"field":"2"}`))
	// THEN
	assert.NoError(t, err)
	assert.Equal(t, int64(0), id)
}

func Test_FileStorage_Append(t *testing.T) {
	// GIVEN
	streamName := "logs"
	storage, err := NewFileStorage(streamName)
	defer func() {
		err = storage.Close()
		cleanFile(t, streamName)
	}()
	_, err = storage.Append(stringToJson(`{"field":"2"}`))
	// WHEN
	var id int64
	id, err = storage.Append(stringToJson(`{"field":"2"}`))
	// THEN
	assert.NoError(t, err)
	assert.Greater(t, id, int64(0))
}

func Test_FileStorage_InitCursor_should_return_0_on_beginning(t *testing.T) {
	// GIVEN
	var storage Storage
	var err error
	streamName := "logs"
	defer func() {
		err = storage.Close()
		cleanFile(t, streamName)
	}()
	storage, err = NewFileStorage(streamName)
	_, err = storage.Append(ast.Json{})
	err = storage.Close()
	storage, err = NewFileStorage(streamName)
	var cursor interface{}
	// WHEN
	cursor, err = storage.InitCursor(model.StartPositionBeginning)
	// THEN
	assert.NoError(t, err)
	assert.Equal(t, 0, cursor)
}

func Test_FileStorage_InitCursor_should_return_21_on_beginning(t *testing.T) {
	// GIVEN
	var storage Storage
	var err error
	streamName := "logs"
	defer func() {
		err = storage.Close()
		cleanFile(t, streamName)
	}()
	storage, err = NewFileStorage(streamName)
	_, err = storage.Append(stringToJson(`{"field":"test"}`))
	err = storage.Close()
	storage, err = NewFileStorage(streamName)
	var cursor interface{}
	// WHEN
	cursor, err = storage.InitCursor(model.StartPositionEnd)
	// THEN
	assert.NoError(t, err)
	assert.Equal(t, int64(21), cursor)
}

func Test_FileStorage_Read_items_1_2(t *testing.T) {
	// GIVEN
	var storage Storage
	var err error
	streamName := "logs"
	defer func() {
		err = storage.Close()
		cleanFile(t, streamName)
	}()
	storage, err = NewFileStorage(streamName)
	for i := 0; i <= 10; i++ {
		_, err = storage.Append(stringToJson(fmt.Sprintf(`{"field":%d}`, i)))
	}
	var cursor interface{}
	cursor, err = storage.InitCursor(model.StartPositionBeginning)
	// WHEN
	cursor, data, endOfStream, err := storage.Read(model.ReadBehaviourNext, cursor, 2)
	// THEN
	assert.NoError(t, err)
	assert.Equal(t, 2, len(data))
	assert.False(t, endOfStream)
}

func Test_FileStorage_Read_items_3_4(t *testing.T) {
	// GIVEN
	var storage Storage
	var err error
	streamName := "logs"
	defer func() {
		err = storage.Close()
		cleanFile(t, streamName)
	}()
	storage, err = NewFileStorage(streamName)
	for i := 0; i <= 10; i++ {
		_, err = storage.Append(stringToJson(fmt.Sprintf(`{"field":%d}`, i)))
	}
	var cursor interface{}
	cursor = int64(32)
	// WHEN
	cursor, data, endOfStream, err := storage.Read(model.ReadBehaviourNext, cursor, 2)
	// THEN
	assert.NoError(t, err)
	assert.Equal(t, 2, len(data))
	assert.False(t, endOfStream)
}

func Test(t *testing.T) {
	assert.Fail(t, "finish read")
}

func cleanFile(t *testing.T, fileName string) {
	if err := os.Remove(fileName); err != nil {
		assert.NoError(t, err)
	}
}
