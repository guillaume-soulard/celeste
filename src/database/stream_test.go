package database

import (
	"celeste/src/model/ast"
	"celeste/src/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Stream_Append_should_append_data_in_down_streams(t *testing.T) {
	// GIVEN
	db := LoadDatabase()
	upStreamName := "upStream"
	err := NewStreamFrom(&db, ast.StreamCreation{
		Name: &upStreamName,
		Storage: &ast.StreamStorage{
			Memory: true,
		},
	})
	assert.NoError(t, err)
	downStreamName := "downStream"
	err = NewStreamFrom(&db, ast.StreamCreation{
		Name: &downStreamName,
		Storage: &ast.StreamStorage{
			Memory: true,
		},
		StreamDataSource: &ast.StreamDataSource{
			From: &upStreamName,
		},
	})
	assert.NoError(t, err)
	data := ast.Json{}
	// WHEN
	id, err := db.Streams[upStreamName].Append(data)
	// THEN
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	upStreamHead := db.Streams[upStreamName].Storage.(*storage.Memory).Data.Head.Data.Data
	downStreamHead := db.Streams[downStreamName].Storage.(*storage.Memory).Data.Head.Data.Data
	assert.Equal(t, upStreamHead, downStreamHead)
}
