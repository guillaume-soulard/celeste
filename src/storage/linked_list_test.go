package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_LinkedList_Append_one_item(t *testing.T) {
	// GIVEN
	list := NewLinkedList[string]()
	// WHEN
	list.Append("test")
	// THEN
	assert.Equal(t, "test", list.Head.Data)
	assert.Equal(t, "test", list.Tail.Data)
	assert.Equal(t, list.Head, list.Tail)
}

func Test_LinkedList_Append_two_items(t *testing.T) {
	// GIVEN
	list := NewLinkedList[string]()
	// WHEN
	list.Append("1")
	list.Append("2")
	// THEN
	assert.Equal(t, "1", list.Head.Data)
	assert.Equal(t, "2", list.Tail.Data)
}

func Test_LinkedList_Append_three_items(t *testing.T) {
	// GIVEN
	list := NewLinkedList[string]()
	// WHEN
	list.Append("1")
	list.Append("2")
	list.Append("3")
	// THEN
	assert.Equal(t, "1", list.Head.Data)
	assert.Equal(t, "2", list.Head.Next.Data)
	assert.Equal(t, "3", list.Tail.Data)
	assert.Equal(t, "2", list.Tail.Previous.Data)
	assert.Equal(t, "1", list.Tail.Previous.Previous.Data)
}
