package ast

type StreamDrop struct {
	Name *string `"DROP" "STREAM" @Ident`
}
