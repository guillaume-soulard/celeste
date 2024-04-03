package storage

import (
	"fmt"
	"time"
)

type IdGenerator struct {
	lastTime time.Time
	sequence int64
}

func NewIdGenerator() IdGenerator {
	return IdGenerator{
		lastTime: time.Now(),
		sequence: 0,
	}
}

func (ig *IdGenerator) NextId(date time.Time) string {
	if ig.lastTime.UnixMilli() == date.UnixMilli() {
		ig.sequence++
	} else {
		ig.sequence = 0
	}
	ig.lastTime = date
	return fmt.Sprintf("%d-%d", ig.lastTime.UnixMilli(), ig.sequence)
}
