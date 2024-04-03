package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_NextId_should_generate_1712265032000_0(t *testing.T) {
	// GIVEN
	idGenerator := NewIdGenerator()
	// WHEN
	id := idGenerator.NextId(time.Date(2024, time.April, 4, 23, 10, 32, 25, time.Local))
	// THEN
	assert.Equal(t, "1712265032000-0", id)
}

func Test_NextId_should_generate_1712265032000_1(t *testing.T) {
	// GIVEN
	idGenerator := NewIdGenerator()
	date := time.Date(2024, time.April, 4, 23, 10, 32, 25, time.Local)
	_ = idGenerator.NextId(date)
	// WHEN
	id := idGenerator.NextId(date)
	// THEN
	assert.Equal(t, "1712265032000-1", id)
}

func Test_NextId_should_generate_1712265032001_0(t *testing.T) {
	// GIVEN
	idGenerator := NewIdGenerator()
	date := time.Date(2024, time.April, 4, 23, 10, 32, 25, time.Local)
	_ = idGenerator.NextId(date)
	// WHEN
	id := idGenerator.NextId(date.Add(time.Millisecond))
	// THEN
	assert.Equal(t, "1712265032001-0", id)
}
