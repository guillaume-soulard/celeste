package storage

import (
	"celeste/src/model"
	"celeste/src/model/ast"
	"errors"
)

func NewMemoryStorage() Storage {
	return &Memory{
		LastId:                0,
		Data:                  NewLinkedList[MemoryData](),
		PreviousReadBehaviour: model.ReadBehaviourNext,
	}
}

type MemoryData struct {
	Id   int64
	Data ast.Json
}

type Memory struct {
	LastId                int64
	Data                  LinkedList[MemoryData]
	PreviousReadBehaviour model.ReadBehaviour
}

func (m *Memory) Append(data ast.Json) (id int64, err error) {
	m.LastId++
	id = m.LastId
	m.Data.Append(MemoryData{
		Id:   m.LastId,
		Data: data,
	})
	return id, err
}

func (m *Memory) InitCursor(startPosition model.StartPosition) (cursor interface{}, err error) {
	if startPosition == model.StartPositionBeginning {
		cursor = m.Data.Head
	} else if startPosition == model.StartPositionEnd {
		cursor = m.Data.Tail
	} else {
		err = errors.New("unsupported start position for storage memory")
	}
	return cursor, err
}

func (m *Memory) Read(readBehaviour model.ReadBehaviour, cursor interface{}, count int) (newCursor interface{}, data []ast.Json, endOfStream bool, err error) {
	data = make([]ast.Json, 0, count)
	node := cursor.(*LinkedListNode[MemoryData])
	for len(data) < count {
		data = append(data, node.Data.Data)
		node = m.getNextNodeFrom(readBehaviour, node)
		if node == nil {
			endOfStream = true
			return newCursor, data, endOfStream, err
		}
	}
	m.PreviousReadBehaviour = readBehaviour
	return newCursor, data, endOfStream, err
}

func (m *Memory) getNextNodeFrom(readBehaviour model.ReadBehaviour, node *LinkedListNode[MemoryData]) *LinkedListNode[MemoryData] {
	if readBehaviour == model.ReadBehaviourNext || (readBehaviour == model.ReadBehaviourAgain && m.PreviousReadBehaviour == model.ReadBehaviourNext) {
		return node.Next
	} else if readBehaviour == model.ReadBehaviourPrevious || (readBehaviour == model.ReadBehaviourAgain && m.PreviousReadBehaviour == model.ReadBehaviourPrevious) {
		return node.Previous
	}
	return nil
}
