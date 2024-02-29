package database

import (
	"celeste/src/model"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDatabase_NewStreamFrom_should_add_stream_in_database(t *testing.T) {
	// GIVEN
	db := LoadDatabase()
	streamName := "foo"
	ast, err := model.Parse(fmt.Sprintf("CREATE STREAM %s STORAGE IN MEMORY", streamName))
	assert.NoError(t, err)
	// WHEN
	err = db.NewStreamFrom(*ast.StreamCreation)
	// THEN
	assert.NoError(t, err)
	assert.Equal(t, 1, len(db.Streams))
	assert.Equal(t, streamName, db.Streams[streamName].Name)
}

func TestDatabase_NewStreamFrom_should_not_add_stream_in_database_because_it_already_exists(t *testing.T) {
	// GIVEN
	db := LoadDatabase()
	streamName := "foo"
	ast, err := model.Parse(fmt.Sprintf("CREATE STREAM %s STORAGE IN MEMORY", streamName))
	assert.NoError(t, err)
	err = db.NewStreamFrom(*ast.StreamCreation)
	// WHEN
	err = db.NewStreamFrom(*ast.StreamCreation)
	// THEN
	assert.Error(t, err)
	assert.Equal(t, 1, len(db.Streams))
	assert.Equal(t, streamName, db.Streams[streamName].Name)
}

func TestDatabase_DropStream_should_delete_stream_in_database(t *testing.T) {
	// GIVEN
	db := LoadDatabase()
	streamName := "foo"
	ast, err := model.Parse(fmt.Sprintf("CREATE STREAM %s STORAGE IN MEMORY", streamName))
	assert.NoError(t, err)
	err = db.NewStreamFrom(*ast.StreamCreation)
	assert.NoError(t, err)
	// WHEN
	err = db.DropStream(streamName)
	// THEN
	assert.NoError(t, err)
	assert.Empty(t, db.Streams)
}

func TestDatabase_DropStream_should_return_error_because_stream_name_not_exists(t *testing.T) {
	// GIVEN
	db := LoadDatabase()
	streamName := "foo"
	// WHEN
	err := db.DropStream(streamName)
	// THEN
	assert.Error(t, err)
	assert.Empty(t, db.Streams)
}
