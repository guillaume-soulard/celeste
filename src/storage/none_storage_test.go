package storage

import (
	"celeste/src/model"
	"celeste/src/model/ast"
	"fmt"
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

func Test_NoneStorage_InitCursor_should_return_0_when_StartPositionBeginning(t *testing.T) {
	// GIVEN
	memory := NewNoneStorage()
	// WHEN
	cursor, err := memory.InitCursor(model.StartPositionBeginning)
	// THEN
	assert.NoError(t, err)
	assert.Equal(t, 0, cursor)
}

func Test_NoneStorage_InitCursor_should_return_0_when_StartPositionEnd(t *testing.T) {
	// GIVEN
	memory := NewNoneStorage()
	// WHEN
	cursor, err := memory.InitCursor(model.StartPositionEnd)
	// THEN
	assert.NoError(t, err)
	assert.Equal(t, 0, cursor)
}

func Test_NoneStorage_Read(t *testing.T) {
	tests := []struct {
		behaviour model.ReadBehaviour
	}{
		{behaviour: model.ReadBehaviourAgain},
		{behaviour: model.ReadBehaviourPrevious},
		{behaviour: model.ReadBehaviourNext},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("should return 0 when read behaviour %d", test.behaviour), func(t *testing.T) {
			// GIVEN
			memory := NewNoneStorage()
			// WHEN
			cursor, data, endOfStream, err := memory.Read(model.ReadBehaviourAgain, 0, 1)
			// THEN
			assert.NoError(t, err)
			assert.Equal(t, 0, cursor)
			assert.Equal(t, []ast.Json{}, data)
			assert.Equal(t, true, endOfStream)
		})
	}
}
