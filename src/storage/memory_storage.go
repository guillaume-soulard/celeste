package storage

import (
	"celeste/src/model"
	"celeste/src/model/ast"
	"errors"
	"time"
	"unsafe"
)

func NewMemoryStorage() Storage {
	return &Memory{
		Size:                  0,
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
	Size                  uint64
	LastId                int64
	Data                  LinkedList[MemoryData]
	PreviousReadBehaviour model.ReadBehaviour
	OldestInsertedTime    time.Time
}

func (m *Memory) Append(data ast.Json) (id int64, err error) {
	m.LastId++
	id = m.LastId
	memoryData := MemoryData{
		Id:   m.LastId,
		Data: data,
	}
	m.Data.Append(memoryData)
	m.Size += uint64(unsafe.Sizeof(memoryData))
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

func (m *Memory) Close() (err error) {
	m.Data.Head = nil
	m.Data.Tail = nil
	return err
}

func (m *Memory) Truncate(evictionPolicies *[]ast.EvictionPolicy) (err error) {
	for _, policy := range *evictionPolicies {
		if policy.MaxAmountItems != nil && uint64(*policy.MaxAmountItems) < m.Data.Len {
			maxLen := uint64(*policy.MaxAmountItems)
			count := uint64(1)
			nodeToDeleteTo := m.Data.Tail
			for nodeToDeleteTo != nil {
				nodeToDeleteTo = nodeToDeleteTo.Previous
				count++
				if count <= maxLen {
					break
				}
			}
			if nodeToDeleteTo != nil {
				nodeToDeleteTo = nodeToDeleteTo.Previous
			}
			for nodeToDeleteTo != nil {
				m.Data.DeleteNode(nodeToDeleteTo)
				nodeToDeleteTo = nodeToDeleteTo.Previous
			}
		}
		if policy.MaxSize != nil && (*policy.MaxSize).Bytes() < m.Size {
			err = errors.New("unsupported max size truncation for storage memory")
		}
		if policy.MaxDuration != nil && time.Now().Sub(m.OldestInsertedTime) > (*policy.MaxDuration).Duration() {
			err = errors.New("unsupported max duration truncation for storage memory")
		}
	}
	return err
}
