package ast

import (
	"strings"
	"time"
)

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

var (
	SizeKb = "KB"
	SizeMb = "MB"
	SizeGb = "GB"
	SizeTb = "TB"
)

var sizeMap = map[string]uint64{
	SizeKb: 1024,
	SizeMb: 1024 * 1024,
	SizeGb: 1024 * 1024 * 1024,
	SizeTb: 1024 * 1024 * 1024 * 1024,
}

func (s Size) Bytes() uint64 {
	if s.Amount == nil {
		return uint64(0)
	}
	if s.Unit == nil {
		return uint64(*s.Amount)
	}
	return uint64(*s.Amount) * sizeMap[*s.Unit]
}

type Duration struct {
	Amount *int    `@Number`
	Unit   *string `@("SECOND" | "SECONDS" | "MINUTE" | "MINUTES" | "HOUR" | "HOURS" | "DAY" | "DAYS" )`
}

var (
	DurationUnitSecond  = "SECOND"
	DurationUnitSeconds = "SECONDS"
	DurationUnitMinute  = "MINUTE"
	DurationUnitMinutes = "MINUTES"
	DurationUnitHour    = "HOUR"
	DurationUnitHours   = "HOURS"
	DurationUnitDay     = "DAY"
	DurationUnitDays    = "DAYS"
)

var durationsByName = map[string]time.Duration{
	strings.ToLower(DurationUnitSecond):  time.Second,
	strings.ToLower(DurationUnitSeconds): time.Second,
	strings.ToLower(DurationUnitMinute):  time.Minute,
	strings.ToLower(DurationUnitMinutes): time.Minute,
	strings.ToLower(DurationUnitHour):    time.Hour,
	strings.ToLower(DurationUnitHours):   time.Hour,
	strings.ToLower(DurationUnitDay):     time.Hour * 24,
	strings.ToLower(DurationUnitDays):    time.Hour * 24,
}

func (d Duration) Duration() time.Duration {
	amount := -1
	unit := time.Second
	if d.Amount != nil && d.Unit != nil {
		amount = *d.Amount
		unit = durationsByName[strings.ToLower(*d.Unit)]
	}
	return unit * time.Duration(amount)
}

type Partition struct {
	ItemAmount *int      `@Number "ITEMS" |`
	Size       *Size     `@@ |`
	Duration   *Duration `@@`
}
