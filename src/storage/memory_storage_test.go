package storage

import (
	"celeste/src/model"
	"celeste/src/model/ast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MemoryStorage_Append_should_add_one_item(t *testing.T) {
	// GIVEN
	memory := NewMemoryStorage()
	data := stringToJson(`{"field":"value"}`)
	// WHEN
	id, err := memory.Append(data)
	// THEN
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
}

func Test_MemoryStorage_Append_should_add_two_items(t *testing.T) {
	// GIVEN
	memory := NewMemoryStorage()
	// WHEN
	_, _ = memory.Append(stringToJson(`{"field":"value"}`))
	id, err := memory.Append(stringToJson(`{"field":"value"}`))
	// THEN
	assert.NoError(t, err)
	assert.Equal(t, int64(2), id)
}

func Test_MemoryStorage_Append_should_add_three_items(t *testing.T) {
	// GIVEN
	memory := NewMemoryStorage().(*Memory)
	// WHEN
	json1 := stringToJson(`{"field":"1"}`)
	json2 := stringToJson(`{"field":"2"}`)
	json3 := stringToJson(`{"field":"3"}`)
	_, _ = memory.Append(json1)
	_, _ = memory.Append(json2)
	id, err := memory.Append(json3)
	// THEN
	assert.NoError(t, err)
	assert.Equal(t, int64(3), id)
	assert.Equal(t, json1, memory.Data.Head.Data.Data)
	assert.Equal(t, json2, memory.Data.Head.Next.Data.Data)
	assert.Equal(t, json3, memory.Data.Head.Next.Next.Data.Data)
}

func Test_MemoryStorage_InitCursor_should_be_set_to_head(t *testing.T) {
	// GIVEN
	memory := NewMemoryStorage().(*Memory)
	_, _ = memory.Append(stringToJson(`{"field":"value"}`))
	_, _ = memory.Append(stringToJson(`{"field":"value"}`))
	_, _ = memory.Append(stringToJson(`{"field":"value"}`))
	// WHEN
	cursor, err := memory.InitCursor(model.StartPositionBeginning)
	// THEN
	assert.NoError(t, err)
	assert.Equal(t, memory.Data.Head, cursor)
}

func Test_MemoryStorage_InitCursor_should_be_set_to_tail(t *testing.T) {
	// GIVEN
	memory := NewMemoryStorage().(*Memory)
	_, _ = memory.Append(stringToJson(`{"field":"value"}`))
	_, _ = memory.Append(stringToJson(`{"field":"value"}`))
	_, _ = memory.Append(stringToJson(`{"field":"value"}`))
	// WHEN
	cursor, err := memory.InitCursor(model.StartPositionEnd)
	// THEN
	assert.NoError(t, err)
	assert.Equal(t, memory.Data.Tail, cursor)
}

func Test_MemoryStorage_InitCursor_should_return_and_error_on_unknown_start_position(t *testing.T) {
	// GIVEN
	memory := NewMemoryStorage().(*Memory)
	// WHEN
	_, err := memory.InitCursor(152)
	// THEN
	assert.Error(t, err)
}

func Test_MemoryStorage_Read_should_return_2_items(t *testing.T) {
	// GIVEN
	memory := NewMemoryStorage().(*Memory)
	json1 := stringToJson(`{"field":"1"}`)
	json2 := stringToJson(`{"field":"2"}`)
	json3 := stringToJson(`{"field":"3"}`)
	_, _ = memory.Append(json1)
	_, _ = memory.Append(json2)
	_, _ = memory.Append(json3)
	cursor, _ := memory.InitCursor(model.StartPositionBeginning)
	// WHEN
	_, result, endOfStream, err := memory.Read(model.ReadBehaviourAgain, cursor, 2)
	// THEN
	assert.Equal(t, 2, len(result))
	assert.Equal(t, json1, result[0])
	assert.Equal(t, json2, result[1])
	assert.NoError(t, err)
	assert.False(t, endOfStream)
}

func Test_MemoryStorage_Read_should_return_3_items(t *testing.T) {
	// GIVEN
	memory := NewMemoryStorage().(*Memory)
	json1 := stringToJson(`{"field":"1"}`)
	json2 := stringToJson(`{"field":"2"}`)
	json3 := stringToJson(`{"field":"3"}`)
	_, _ = memory.Append(json1)
	_, _ = memory.Append(json2)
	_, _ = memory.Append(json3)
	cursor, _ := memory.InitCursor(model.StartPositionBeginning)
	// WHEN
	_, result, endOfStream, err := memory.Read(model.ReadBehaviourAgain, cursor, 3)
	// THEN
	assert.Equal(t, 3, len(result))
	if len(result) == 3 {
		assert.Equal(t, json1, result[0])
		assert.Equal(t, json2, result[1])
		assert.Equal(t, json3, result[2])
	}
	assert.NoError(t, err)
	assert.True(t, endOfStream)
}

func Test_MemoryStorage_Read_should_return_3_items_with_count_5(t *testing.T) {
	// GIVEN
	memory := NewMemoryStorage().(*Memory)
	json1 := stringToJson(`{"field":"1"}`)
	json2 := stringToJson(`{"field":"2"}`)
	json3 := stringToJson(`{"field":"3"}`)
	_, _ = memory.Append(json1)
	_, _ = memory.Append(json2)
	_, _ = memory.Append(json3)
	cursor, _ := memory.InitCursor(model.StartPositionBeginning)
	// WHEN
	_, result, endOfStream, err := memory.Read(model.ReadBehaviourNext, cursor, 5)
	// THEN
	assert.Equal(t, 3, len(result))
	if len(result) == 3 {
		assert.Equal(t, json1, result[0])
		assert.Equal(t, json2, result[1])
		assert.Equal(t, json3, result[2])
	}
	assert.NoError(t, err)
	assert.True(t, endOfStream)
}

func Test_MemoryStorage_Read_should_return_2_items_from_end(t *testing.T) {
	// GIVEN
	memory := NewMemoryStorage().(*Memory)
	json1 := stringToJson(`{"field":"1"}`)
	json2 := stringToJson(`{"field":"2"}`)
	json3 := stringToJson(`{"field":"3"}`)
	_, _ = memory.Append(json1)
	_, _ = memory.Append(json2)
	_, _ = memory.Append(json3)
	cursor, _ := memory.InitCursor(model.StartPositionEnd)
	// WHEN
	_, result, endOfStream, err := memory.Read(model.ReadBehaviourPrevious, cursor, 2)
	// THEN
	assert.Equal(t, 2, len(result))
	if len(result) == 2 {
		assert.Equal(t, json3, result[0])
		assert.Equal(t, json2, result[1])
	}
	assert.NoError(t, err)
	assert.False(t, endOfStream)
}

var amount2 = 2

func Test_MemoryStorage_Truncate(t *testing.T) {
	// GIVEN
	memory := NewMemoryStorage().(*Memory)
	json1 := stringToJson(`{"field":"1"}`)
	json2 := stringToJson(`{"field":"2"}`)
	json3 := stringToJson(`{"field":"3"}`)
	_, _ = memory.Append(json1)
	_, _ = memory.Append(json2)
	_, _ = memory.Append(json3)
	// WHEN
	err := memory.Truncate(&[]ast.EvictionPolicy{
		{
			MaxAmountItems: &amount2,
		},
	})
	// THEN
	assert.Equal(t, uint64(2), memory.Len)
	assert.NoError(t, err)
}

func stringToJson(str string) ast.Json {
	parsed, _ := model.Parse(str)
	return *parsed.Json
}
