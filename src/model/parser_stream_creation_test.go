package model

import (
	"testing"
)

func Test_stream_creation(t *testing.T) {
	tests := []commandTest{
		{
			Command:       `CREATE STREAM 123`,
			ErrorExpected: true,
		},
		{
			Command:       `CREATE STREAM`,
			ErrorExpected: true,
		},
		{
			Command:       `CREATE STREAM test`,
			ErrorExpected: false,
		},
		{
			Command:       `create database TEST`,
			ErrorExpected: false,
		},
		{
			Command:       `CREATE STREAM test unknown`,
			ErrorExpected: true,
		},
		{
			Command:       `CREATE STREAM test STORAGE IN MEMORY`,
			ErrorExpected: false,
		},
		{
			Command:       `CREATE STREAM test STORAGE IN MEMORY FROM test`,
			ErrorExpected: false,
		},
		{
			Command:       `CREATE STREAM test STORAGE IN MEMORY FROM test FILTER 'test'`,
			ErrorExpected: false,
		},
		{
			Command:       `CREATE STREAM test STORAGE IN MEMORY FROM test FILTER 'test' MAP 'test'`,
			ErrorExpected: false,
		},
		{
			Command:       `CREATE STREAM test STORAGE IN MEMORY FROM test FILTER 'test' MAP 'test' EVICTION POLICY MAX SIZE 10 MB`,
			ErrorExpected: false,
		},
		{
			Command:       `CREATE STREAM test STORAGE IN MEMORY FROM test FILTER 'test' MAP 'test' EVICTION POLICY MAX SIZE 10 MB AND MAX ITEMS 15`,
			ErrorExpected: false,
		},
		{
			Command:       `CREATE STREAM test STORAGE IN MEMORY FROM test FILTER 'test' MAP 'test' EVICTION POLICY MAX SIZE 10 MB PARTITION BY 10 MINUTES`,
			ErrorExpected: false,
		},
	}
	runTest(t, tests)
}
