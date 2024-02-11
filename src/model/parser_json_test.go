package model

import (
	"testing"
)

func Test_stream_json(t *testing.T) {
	tests := []commandTest{
		{
			Command:       `{}`,
			ErrorExpected: false,
		},
		{
			Command:       `{"field": "value"}`,
			ErrorExpected: false,
		},
		{
			Command:       `{"field": "value","field2":1}`,
			ErrorExpected: false,
		},
		{
			Command:       `{"field": 1}`,
			ErrorExpected: false,
		},
		{
			Command:       `{"field": true}`,
			ErrorExpected: false,
		},
		{
			Command:       `{"field": false}`,
			ErrorExpected: false,
		},
		{
			Command:       `{"field": null}`,
			ErrorExpected: false,
		},
		{
			Command:       `[{"field": null}]`,
			ErrorExpected: false,
		},
		{
			Command:       `{"field": { "subField": true }}`,
			ErrorExpected: false,
		},
		{
			Command:       `null`,
			ErrorExpected: false,
		},
		{
			Command:       `1`,
			ErrorExpected: false,
		},
		{
			Command:       `"test"`,
			ErrorExpected: false,
		},
		{
			Command:       `true`,
			ErrorExpected: false,
		},
		{
			Command:       `false`,
			ErrorExpected: false,
		},
		{
			Command:       `[]`,
			ErrorExpected: false,
		},
		{
			Command:       `[null, true, false, 1, "test"]`,
			ErrorExpected: false,
		},
	}
	runTest(t, tests)
}
