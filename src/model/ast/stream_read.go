package ast

type StreamRead struct {
	Stop          bool           `"READ" (@"STOP" | `
	ReadDirection *ReadDirection `@@ |`
	ReadClause    *ReadClause    `"START" @@ )`
}

type ReadDirection struct {
	Previous bool `@"PREVIOUS" |`
	Again    bool `@"AGAIN" |`
	Next     bool `@"NEXT"`
}

type ReadClause struct {
	Stream           *string      `@Ident`
	ConsumerName     *string      `["WITH" "CONSUMER" @Ident]`
	StartAtBeginning bool         `["START" "AT" (@"BEGINNING" |`
	StartAtEnd       bool         `@"END")]`
	Predicates       *[]Predicate `@@*`
	Count            *int         `["COUNT" @Number]`
	Follow           bool         `[@"FOLLOW"]`
	WrireTo          []string     `["WRITE" "TO" @Ident ("," @Ident)*]`
}
