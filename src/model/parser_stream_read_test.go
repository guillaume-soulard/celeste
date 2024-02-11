package model

import (
	"testing"
)

func Test_stream_read(t *testing.T) {
	tests := []commandTest{
		{
			Command:       `read`,
			ErrorExpected: true,
		},
		{
			Command:       `read stop`,
			ErrorExpected: false,
		},
		{
			Command:       `read stop test`,
			ErrorExpected: true,
		},
		{
			Command:       `read previous`,
			ErrorExpected: false,
		},
		{
			Command:       `read previous test`,
			ErrorExpected: true,
		},
		{
			Command:       `read again`,
			ErrorExpected: false,
		},
		{
			Command:       `read again test`,
			ErrorExpected: true,
		},
		{
			Command:       `read next`,
			ErrorExpected: false,
		},
		{
			Command:       `read next test`,
			ErrorExpected: true,
		},
	}
	runTest(t, tests)
}

func Test_stream_read_start(t *testing.T) {
	tests := []commandTest{
		{
			Command:       `read start`,
			ErrorExpected: true,
		},
		{
			Command:       `read start test`,
			ErrorExpected: false,
		},
		{
			Command:       `read start test with consumer test`,
			ErrorExpected: false,
		},
		{
			Command:       `read start test with`,
			ErrorExpected: true,
		},
		{
			Command:       `read start test with consumer`,
			ErrorExpected: true,
		},
		{
			Command:       `read start test start`,
			ErrorExpected: true,
		},
		{
			Command:       `read start test start at`,
			ErrorExpected: true,
		},
		{
			Command:       `read start test start at test`,
			ErrorExpected: true,
		},
		{
			Command:       `read start test start at beginning`,
			ErrorExpected: false,
		},
		{
			Command:       `read start test start at end`,
			ErrorExpected: false,
		},
		{
			Command:       `read start test count`,
			ErrorExpected: true,
		},
		{
			Command:       `read start test count test`,
			ErrorExpected: true,
		},
		{
			Command:       `read start test count 10`,
			ErrorExpected: false,
		},
		{
			Command:       `read start test follow`,
			ErrorExpected: false,
		},
		{
			Command:       `read start test write`,
			ErrorExpected: true,
		},
		{
			Command:       `read start test write to`,
			ErrorExpected: true,
		},
		{
			Command:       `read start test write to 10`,
			ErrorExpected: true,
		},
		{
			Command:       `read start test write test`,
			ErrorExpected: true,
		},
		{
			Command:       `read start test write to test`,
			ErrorExpected: false,
		},
		{
			Command:       `read start test filter`,
			ErrorExpected: true,
		},
		{
			Command:       `read start test filter test`,
			ErrorExpected: true,
		},
		{
			Command:       `read start test filter 'test'`,
			ErrorExpected: false,
		},
		{
			Command:       `read start test map`,
			ErrorExpected: true,
		},
		{
			Command:       `read start test map test`,
			ErrorExpected: true,
		},
		{
			Command:       `read start test map 'test'`,
			ErrorExpected: false,
		},
		{
			Command: `read start test
with consumer consumer
start at beginning
filter 'field = value'
map 'field'
count 10
follow
write to stream1, stream2
`,
			ErrorExpected: false,
		},
	}
	runTest(t, tests)
}
