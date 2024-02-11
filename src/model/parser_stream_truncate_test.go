package model

import (
	"testing"
)

func Test_stream_truncate(t *testing.T) {
	tests := []commandTest{
		{
			Command:       `TRUNCATE STREAM 123`,
			ErrorExpected: true,
		},
		{
			Command:       `TRUNCATE STREAM`,
			ErrorExpected: true,
		},
		{
			Command:       `TRUNCATE STREAM test unknown`,
			ErrorExpected: true,
		},
		{
			Command:       `TRUNCATE STREAM test`,
			ErrorExpected: false,
		},
		{
			Command:       `TRUNCATE STREAM test WITH MAX ITEMS 10`,
			ErrorExpected: false,
		},
		{
			Command:       `TRUNCATE STREAM test WITH MAX SIZE 15 MB`,
			ErrorExpected: false,
		},
		{
			Command:       `TRUNCATE STREAM test WITH MAX DURATION 10 SECONDS`,
			ErrorExpected: false,
		},
		{
			Command:       `TRUNCATE STREAM test WITH MAX DURATION 10 SECONDS AND MAX SIZE 15 MB`,
			ErrorExpected: false,
		},
	}
	runTest(t, tests)
}
