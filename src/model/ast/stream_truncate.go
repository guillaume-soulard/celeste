package ast

type StreamTruncate struct {
	Name             *string           `"TRUNCATE" "STREAM" @Ident`
	EvictionPolicies *[]EvictionPolicy `["WITH" @@ ("AND" @@)*]`
}
