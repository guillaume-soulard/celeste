package ast

type StreamCreation struct {
	Name             *string           `"CREATE" "STREAM" @Ident`
	Storage          *StreamStorage    `["STORAGE" @@]`
	StreamDataSource *StreamDataSource `[ @@ ]`
	EvictionPolicy   *[]EvictionPolicy `["EVICTION" "POLICY" @@ ("AND" @@)*]`
	Partition        *Partition        `["PARTITION" "BY" @@]`
}

type StreamStorage struct {
	None   bool `@"NO" |`
	Memory bool `@("IN" "MEMORY") |`
	Disk   bool `@("ON" "DISK")`
}

type StreamDataSource struct {
	From       *string      `"FROM" @Ident`
	Predicates *[]Predicate `@@*`
}

type Predicate struct {
	Filter *Filter `@@ |`
	Map    *Map    `@@`
}

type Filter struct {
	Expression *string `"FILTER" @String`
}

type Map struct {
	Expression *string `"MAP" @String`
}

type EvictionPolicy struct {
	MaxAmountItems *int      `"MAX" "ITEMS" @Number |`
	MaxSize        *Size     `"MAX" "SIZE" @@ |`
	MaxDuration    *Duration `"MAX" "DURATION" @@`
}

type Size struct {
	Amount *int    `@Number`
	Unit   *string `@("KB" | "MB" | "GB" | "TB")`
}

type Duration struct {
	Amount *int    `@Number`
	Unit   *string `@("SECOND" | "SECONDS" | "MINUTE" | "MINUTES" | "HOUR" | "HOURS" | "DAY" | "DAYS" | "MONTH" | "MONTHS" | "YEAR" | "YEARS" )`
}

type Partition struct {
	ItemAmount *int      `@Number "ITEMS" |`
	Size       *Size     `@@ |`
	Duration   *Duration `@@`
}
