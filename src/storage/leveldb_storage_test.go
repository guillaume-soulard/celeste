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

func Test_LevelDbStorage_Read_should_return_items(t *testing.T) {
	tests := []struct {
		numberOfItems       int
		startPosition       model.StartPosition
		readBehaviour       model.ReadBehaviour
		readCount           int
		expectedCount       int
		expectedEndOfStream bool
	}{
		{
			numberOfItems:       10,
			startPosition:       model.StartPositionBeginning,
			readBehaviour:       model.ReadBehaviourNext,
			readCount:           10,
			expectedCount:       10,
			expectedEndOfStream: true,
		},
		{
			numberOfItems:       11,
			startPosition:       model.StartPositionBeginning,
			readBehaviour:       model.ReadBehaviourNext,
			readCount:           10,
			expectedCount:       10,
			expectedEndOfStream: false,
		},
		{
			numberOfItems:       0,
			startPosition:       model.StartPositionBeginning,
			readBehaviour:       model.ReadBehaviourNext,
			readCount:           10,
			expectedCount:       0,
			expectedEndOfStream: true,
		},
		{
			numberOfItems:       5,
			startPosition:       model.StartPositionBeginning,
			readBehaviour:       model.ReadBehaviourNext,
			readCount:           10,
			expectedCount:       5,
			expectedEndOfStream: true,
		},
		{
			numberOfItems:       5,
			startPosition:       model.StartPositionEnd,
			readBehaviour:       model.ReadBehaviourPrevious,
			readCount:           2,
			expectedCount:       2,
			expectedEndOfStream: false,
		},
	}
	for _, test := range tests {
		name := fmt.Sprintf(
			"read should read %d items and end of stream %v from db of %d items, read behaviour %v start position %v and read count %d",
			test.expectedCount,
			test.expectedEndOfStream,
			test.numberOfItems,
			test.readBehaviour,
			test.startPosition,
			test.readCount,
		)
		t.Run(name, func(t *testing.T) {
			// GIVEN
			streamName := "logs"
			storage, err := NewLevelDbStorage(streamName)
			defer func() {
				err = storage.Close()
				cleanDir(t, streamName)
			}()
			// WHEN
			for i := 0; i < test.numberOfItems; i++ {
				_, err = storage.Append(stringToJson(fmt.Sprintf(`{"field":"%d"}`, i)))
			}
			assert.NoError(t, err)
			var cursor interface{}
			cursor, err = storage.InitCursor(test.startPosition)
			var data []ast.Json
			var endOfStream bool
			cursor, data, endOfStream, err = storage.Read(test.readBehaviour, cursor, test.readCount)
			// THEN
			assert.Equal(t, test.expectedEndOfStream, endOfStream)
			assert.Equal(t, test.expectedCount, len(data))
			assert.NoError(t, err)
		})
	}
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
