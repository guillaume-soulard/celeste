package ast

type StreamTruncate struct {
	Name           *string           `"TRUNCATE" "STREAM" @Ident`
	EvictionPolicy *[]EvictionPolicy `["WITH" @@ ("AND" @@)*]`
}
