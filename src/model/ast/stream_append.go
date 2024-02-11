package ast

type StreamAppend struct {
	StreamNames *[]string `@Ident ("," @Ident)* "<"`
	JsonData    *[]Json   `@@ ("," @@)*`
}
