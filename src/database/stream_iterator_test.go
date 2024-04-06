package database

import (
	"celeste/src/model/ast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Read_should_return_error_when_count_is_negative(t *testing.T) {
	// GIVEN
	si := StreamIterator{
		Count: -1,
	}
	// WHEN
	err := si.Read(&ast.StreamRead{})
	// THEN
	assert.Error(t, err)
}
