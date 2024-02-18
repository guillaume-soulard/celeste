package storage

import (
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
	assert.Equal(t, int64(0), id)
}

func cleanFile(t *testing.T, fileName string) {
	if err := os.Remove(fileName); err != nil {
		assert.NoError(t, err)
	}
}
