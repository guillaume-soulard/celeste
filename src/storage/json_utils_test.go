package storage

import (
	"celeste/src/model/ast"
	"encoding/json"
	"github.com/oliveagle/jsonpath"
	"github.com/stretchr/testify/assert"
	"testing"
)

func JsonPathLookUp(t *testing.T, j *ast.Json, path string) interface{} {
	jsonString := j.ToString()
	var obj interface{}
	err := json.Unmarshal([]byte(jsonString), &obj)
	assert.NoError(t, err)
	result, err := jsonpath.JsonPathLookup(obj, path)
	assert.NoError(t, err)
	return result
}
