package storage

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_FileStorage_NewFileStorage_should_return_storage(t *testing.T) {
	// GIVEN
	defer func() {
		cleanFile(t)
	}()
	streamName := "logs"
	// WHEN
	storage, err := NewFileStorage(streamName)
	// THEN
	assert.NoError(t, err)
	assert.NotNil(t, storage)
}

func cleanFile(t *testing.T) {
	if err := os.Remove("logs"); err != nil {
		assert.NoError(t, err)
	}
}
