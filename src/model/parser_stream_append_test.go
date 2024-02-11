package model

import (
	"testing"
)

func Test_stream_append(t *testing.T) {
	tests := []commandTest{
		{
			Command:       `test`,
			ErrorExpected: true,
		},
		{
			Command:       `test <`,
			ErrorExpected: true,
		},
		{
			Command:       `test, test <`,
			ErrorExpected: true,
		},
		{
			Command:       `test, test < {"field": "value"}`,
			ErrorExpected: false,
		},
		{
			Command:       `test, test < {}, {}`,
			ErrorExpected: false,
		},
		{
			Command:       `test, test < []`,
			ErrorExpected: false,
		},
		{
			Command:       `test, test < [{"field":"string","field2":2}, true, 1]`,
			ErrorExpected: false,
		},
		{
			Command:       `test, test < null`,
			ErrorExpected: false,
		},
		{
			Command:       `test, test < "test"`,
			ErrorExpected: false,
		},
	}
	runTest(t, tests)
}
