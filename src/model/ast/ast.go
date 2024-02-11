package ast

type Ast struct {
	Json           *Json           `@@ |`
	StreamCreation *StreamCreation `@@ |`
	StreamDrop     *StreamDrop     `@@ |`
	StreamAppend   *StreamAppend   `@@ |`
	StreamTruncate *StreamTruncate `@@ |`
	StreamRead     *StreamRead     `@@`
}
