package ast

import (
	"fmt"
	"strings"
)

type Json struct {
	Value Value `@@`
}

func (j *Json) ToString() string {
	return j.Value.ToString()
}

type JsonObject struct {
	Fields []JsonField `"{" [@@]("," @@)* "}"`
}

type JsonField struct {
	Name  string `@JsonString ":"`
	Value Value  `@@`
}

type Value struct {
	String     *string     `(@JsonString |`
	Number     *float64    `@Number |`
	Boolean    *Boolean    `@@ |`
	Null       bool        `@"NULL" |`
	JsonObject *JsonObject `@@ |`
	JsonArray  *JsonArray  `@@)`
}

type Boolean struct {
	True  bool `@"TRUE" |`
	False bool `@"FALSE"`
}

func (v *Value) ToString() string {
	if v.Null {
		return "null"
	} else if v.String != nil {
		return fmt.Sprintf(`"%s"`, *v.String)
	} else if v.Boolean != nil {
		if v.Boolean.True {
			return "true"
		} else if v.Boolean.False {
			return "false"
		}
	} else if v.Number != nil {
		return fmt.Sprintf("%v", *v.Number)
	} else if v.JsonObject != nil {
		fieldStrs := make([]string, len(v.JsonObject.Fields))
		for i, field := range v.JsonObject.Fields {
			fieldStrs[i] = fmt.Sprintf(`"%s": %v`, field.Name, field.Value.ToString())
		}
		return fmt.Sprintf("{%s}", strings.Join(fieldStrs, ","))
	} else if v.JsonArray != nil {
		arrayLen := 0
		if v.JsonArray.JsonArray != nil {
			arrayLen = len(*v.JsonArray.JsonArray)
		}
		items := make([]string, arrayLen)
		if v.JsonArray.JsonArray != nil {
			for i, item := range *v.JsonArray.JsonArray {
				items[i] = fmt.Sprintf(`%v`, item.ToString())
			}
		}
		return fmt.Sprintf("[%s]", strings.Join(items, ","))
	}
	return ""
}

type JsonArray struct {
	JsonArray *[]Value `"[" [@@] ("," @@)* "]"`
}
