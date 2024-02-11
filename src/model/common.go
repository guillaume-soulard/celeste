package model

type ReadBehaviour byte

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
