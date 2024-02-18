package model

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_Json_ToString(t *testing.T) {
	tests := []struct {
		json string
	}{
		{
			json: `{}`,
		},
		{
			json: `{"field": "value"}`,
		},
		{
			json: `{"field": "value","field2":1}`,
		},
		{
			json: `{"field": 1}`,
		},
		{
			json: `{"field": true}`,
		},
		{
			json: `{"field": false}`,
		},
		{
			json: `{"field": null}`,
		},
		{
			json: `[{"field": null}]`,
		},
		{
			json: `{"field": { "subField": true }}`,
		},
		{
			json: `null`,
		},
		{
			json: `1`,
		},
		{
			json: `"test"`,
		},
		{
			json: `true`,
		},
		{
			json: `false`,
		},
		{
			json: `[]`,
		},
		{
			json: `[null, true, false, 1, "test"]`,
		},
	}
	for _, test := range tests {
		t.Run(test.json, func(t *testing.T) {
			// GIVEN
			fmt.Println(fmt.Sprintf("Testing %s", test.json))
			parsed, err := Parse(test.json)
			assert.NoError(t, err)
			assert.NotNil(t, parsed.Json)
			parsedJson := *parsed.Json
			// WHEN
			result := parsedJson.ToString()
			// THEN
			assert.Equal(t, strings.ReplaceAll(test.json, " ", ""), strings.ReplaceAll(result, " ", ""))
		})
	}
}
