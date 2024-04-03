package model

import (
	"fmt"
	"math"
)

type ReadBehaviour byte

var (
	MinId = "0-0"
	MaxId = fmt.Sprintf("%d-%d", math.MaxInt64, math.MaxInt64)
)

const (
	ReadBehaviourAgain    ReadBehaviour = 0
	ReadBehaviourNext     ReadBehaviour = 1
	ReadBehaviourPrevious ReadBehaviour = 2
)

type StartPosition byte

const (
	StartPositionBeginning StartPosition = 0
	StartPositionEnd       StartPosition = 1
)
