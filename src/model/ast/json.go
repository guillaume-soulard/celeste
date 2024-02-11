package ast

type Json struct {
	Value Value `@@`
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
	Boolean    *bool       `@("TRUE" | "FALSE") |`
	Null       bool        `@"NULL" |`
	JsonObject *JsonObject `@@ |`
	JsonArray  *JsonArray  `@@)`
}

type JsonArray struct {
	JsonArray *[]Value `"[" [@@] ("," @@)* "]"`
}
