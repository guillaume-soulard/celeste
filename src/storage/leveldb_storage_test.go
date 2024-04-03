package storage

import (
	"celeste/src/model"
	"celeste/src/model/ast"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_LevelDbStorage_NewLevelDbStorage_should_return_storage(t *testing.T) {
	// GIVEN
	var storage Storage
	var err error
	streamName := "logs"
	defer func() {
		err = storage.Close()
		cleanDir(t, streamName)
	}()
	// WHEN
	storage, err = NewLevelDbStorage(streamName)
	// THEN
	assert.NoError(t, err)
	assert.NotNil(t, storage)
}

func Test_LevelDbStorage_NewLevelDbStorage_should_append_read_from_beginning(t *testing.T) {
	// GIVEN
	streamName := "logs"
	storage, err := NewLevelDbStorage(streamName)
	defer func() {
		err = storage.Close()
		cleanDir(t, streamName)
	}()
	var id string
	// WHEN
	for i := 0; i < 100; i++ {
		id, err = storage.Append(stringToJson(fmt.Sprintf(`{"field":"%d"}`, i)))
	}
	assert.NoError(t, err)
	var cursor interface{}
	cursor, err = storage.InitCursor(model.StartPositionBeginning)
	var data []ast.Json
	var endOfStream bool
	cursor, data, endOfStream, err = storage.Read(model.ReadBehaviourNext, cursor, 10)
	// THEN
	assert.Equal(t, false, endOfStream)
	assert.Equal(t, 10, len(data))
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
}

func Test_LevelDbStorage_NewLevelDbStorage_should_read_from_end(t *testing.T) {
	// GIVEN
	streamName := "logs"
	storage, err := NewLevelDbStorage(streamName)
	defer func() {
		err = storage.Close()
		cleanDir(t, streamName)
	}()
	var id string
	// WHEN
	for i := 0; i < 100; i++ {
		id, err = storage.Append(stringToJson(fmt.Sprintf(`{"field":"%d"}`, i)))
	}
	assert.NoError(t, err)
	var cursor interface{}
	cursor, err = storage.InitCursor(model.StartPositionEnd)
	var data []ast.Json
	var endOfStream bool
	cursor, data, endOfStream, err = storage.Read(model.ReadBehaviourPrevious, cursor, 10)
	// THEN
	assert.Equal(t, false, endOfStream)
	assert.Equal(t, 10, len(data))
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
}

func cleanDir(t *testing.T, fileName string) {
	var err error
	if err = os.RemoveAll(fileName); err != nil {
		assert.NoError(t, err)
	}
}
