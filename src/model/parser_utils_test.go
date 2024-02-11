package model

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type commandTest struct {
	Command       string
	ErrorExpected bool
}

func runTest(t *testing.T, tests []commandTest) {
	for _, test := range tests {
		testName := fmt.Sprintf("command %s should return error %v", test.Command, test.ErrorExpected)
		t.Run(testName, func(t *testing.T) {
			// GIVEN
			// WHEN
			ast, err := Parse(test.Command)
			// THEN
			if test.ErrorExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, ast)
			}
		})
	}
}
