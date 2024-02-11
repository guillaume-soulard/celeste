package model

import (
	"testing"
)

func Test_stream_drop(t *testing.T) {
	tests := []commandTest{
		{
			Command:       `DROP STREAM 123`,
			ErrorExpected: true,
		},
		{
			Command:       `DROP STREAM`,
			ErrorExpected: true,
		},
		{
			Command:       `DROP STREAM test unknown`,
			ErrorExpected: true,
		},
		{
			Command:       `DROP STREAM test`,
			ErrorExpected: false,
		},
	}
	runTest(t, tests)
}
