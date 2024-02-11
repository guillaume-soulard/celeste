package storage

import (
	"celeste/src/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NoneStorage_Append_when_one_item_is_added(t *testing.T) {
	// GIVEN
	memory := NewNoneStorage()
	parsed, _ := model.Parse(`{"field":"value"}`)
	data := *parsed.Json
	// WHEN
	id, err := memory.Append(data)
	// THEN
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
}

func Test_NoneStorage_Append_when_two_items_is_added(t *testing.T) {
	// GIVEN
	memory := NewNoneStorage()
	parsed, _ := model.Parse(`{"field":"value"}`)
	data := *parsed.Json
	// WHEN
	_, _ = memory.Append(data)
	id, err := memory.Append(data)
	// THEN
	assert.NoError(t, err)
	assert.Equal(t, int64(2), id)
}
