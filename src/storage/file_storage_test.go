package storage

import (
	"celeste/src/model"
	"celeste/src/model/ast"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
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
	var id string
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
	var id string
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
	assert.Equal(t, int64(0), cursor)
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
	assert.Equal(t, int64(25), cursor)
}

type jsonPathTest struct {
	jsonPath      string
	expectedValue interface{}
}

type readTest struct {
	nbOfItemsToTest     int
	startPositionToTest model.StartPosition
	readBehaviourToTest model.ReadBehaviour
	cursor              int64
	readCount           int
	expectedJsonPath    []jsonPathTest
	expectedEndOfStream bool
}

func Test_FileStorage_Read(t *testing.T) {
	tests := []readTest{
		{
			nbOfItemsToTest:     0,
			startPositionToTest: model.StartPositionBeginning,
			readBehaviourToTest: model.ReadBehaviourNext,
			cursor:              -1,
			readCount:           2,
			expectedJsonPath:    []jsonPathTest{},
			expectedEndOfStream: true,
		},
		{
			nbOfItemsToTest:     0,
			startPositionToTest: model.StartPositionBeginning,
			readBehaviourToTest: model.ReadBehaviourAgain,
			cursor:              -1,
			readCount:           2,
			expectedJsonPath:    []jsonPathTest{},
			expectedEndOfStream: true,
		},
		{
			nbOfItemsToTest:     0,
			startPositionToTest: model.StartPositionBeginning,
			readBehaviourToTest: model.ReadBehaviourPrevious,
			cursor:              -1,
			readCount:           2,
			expectedJsonPath:    []jsonPathTest{},
			expectedEndOfStream: true,
		},
		{
			nbOfItemsToTest:     0,
			startPositionToTest: model.StartPositionEnd,
			readBehaviourToTest: model.ReadBehaviourNext,
			cursor:              -1,
			readCount:           2,
			expectedJsonPath:    []jsonPathTest{},
			expectedEndOfStream: true,
		},
		{
			nbOfItemsToTest:     10,
			startPositionToTest: model.StartPositionBeginning,
			readBehaviourToTest: model.ReadBehaviourNext,
			cursor:              -1,
			readCount:           2,
			expectedJsonPath: []jsonPathTest{
				{jsonPath: "$.field", expectedValue: float64(1)},
				{jsonPath: "$.field", expectedValue: float64(2)},
			},
			expectedEndOfStream: false,
		},
		{
			nbOfItemsToTest:     10,
			startPositionToTest: model.StartPositionBeginning,
			readBehaviourToTest: model.ReadBehaviourNext,
			cursor:              int64(40),
			readCount:           2,
			expectedJsonPath: []jsonPathTest{
				{jsonPath: "$.field", expectedValue: float64(3)},
				{jsonPath: "$.field", expectedValue: float64(4)},
			},
			expectedEndOfStream: false,
		},
		{
			nbOfItemsToTest:     3,
			startPositionToTest: model.StartPositionBeginning,
			readBehaviourToTest: model.ReadBehaviourNext,
			cursor:              int64(40),
			readCount:           2,
			expectedJsonPath: []jsonPathTest{
				{jsonPath: "$.field", expectedValue: float64(3)},
			},
			expectedEndOfStream: true,
		},
		{
			nbOfItemsToTest:     10,
			startPositionToTest: model.StartPositionEnd,
			readBehaviourToTest: model.ReadBehaviourNext,
			cursor:              -1,
			readCount:           2,
			expectedJsonPath:    []jsonPathTest{},
			expectedEndOfStream: true,
		},
		{
			nbOfItemsToTest:     10,
			startPositionToTest: model.StartPositionEnd,
			readBehaviourToTest: model.ReadBehaviourPrevious,
			cursor:              -1,
			readCount:           2,
			expectedJsonPath: []jsonPathTest{
				{jsonPath: "$.field", expectedValue: float64(10)},
				{jsonPath: "$.field", expectedValue: float64(9)},
			},
			expectedEndOfStream: false,
		},
	}
	for testIndex, test := range tests {
		paths := make([]string, len(test.expectedJsonPath))
		for i, path := range test.expectedJsonPath {
			paths[i] = fmt.Sprintf("%s = %v", path.jsonPath, path.expectedValue)
		}
		t.Run(fmt.Sprintf("test index %d : should return %d items with paths %s", testIndex, len(test.expectedJsonPath), strings.Join(paths, " and ")), func(t *testing.T) {
			// GIVEN
			var storage Storage
			var err error
			streamName := "logs"
			cleanFile(t, streamName)
			defer func() {
				err = storage.Close()
				cleanFile(t, streamName)
			}()
			storage, err = NewFileStorage(streamName)
			for i := 1; i <= test.nbOfItemsToTest; i++ {
				_, err = storage.Append(stringToJson(fmt.Sprintf(`{"field":%d}`, i)))
			}
			var cursor interface{}
			if test.cursor < 0 {
				cursor, err = storage.InitCursor(test.startPositionToTest)
			} else {
				cursor = test.cursor
			}
			// WHEN
			newCursor, data, endOfStream, err := storage.Read(test.readBehaviourToTest, cursor, test.readCount)
			// THEN
			assert.NotNil(t, newCursor)
			assert.NoError(t, err)
			assert.Equal(t, len(test.expectedJsonPath), len(data))
			for i, expectedPath := range test.expectedJsonPath {
				assert.Equal(t, expectedPath.expectedValue, JsonPathLookUp(t, &data[i], expectedPath.jsonPath))
			}
			assert.Equal(t, test.expectedEndOfStream, endOfStream)
		})
	}
}

func Test(t *testing.T) {
	// GIVEN
	var storage Storage
	var err error
	streamName := "logs"
	cleanFile(t, streamName)
	defer func() {
		err = storage.Close()
		cleanFile(t, streamName)
	}()
	storage, err = NewFileStorage(streamName)
	for i := 1; i <= 10; i++ {
		_, err = storage.Append(stringToJson(fmt.Sprintf(`{"field":%d}`, i)))
	}
	// WHEN
	err = storage.Truncate(&[]ast.EvictionPolicy{
		{
			MaxAmountItems: nil,
			MaxSize:        nil,
			MaxDuration:    nil,
		},
	})
	// THEN
	assert.NoError(t, err)
}

func cleanFile(t *testing.T, fileName string) {
	var err error
	var file *os.File
	if file, err = os.Open(fileName); os.IsNotExist(err) {
		return
	}
	if err = file.Close(); err != nil {
		assert.NoError(t, err)
	}
	if err = os.Remove(fileName); err != nil {
		assert.NoError(t, err)
	}
}
