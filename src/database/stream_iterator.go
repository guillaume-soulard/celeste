package database

import (
	"celeste/src/model"
	"celeste/src/model/ast"
	"fmt"
)

type StreamIterator struct {
	Stream        *Stream
	Cursor        interface{}
	StartPosition model.StartPosition
	Count         int
}

func (i *StreamIterator) Read(db *Database, read *ast.StreamRead) (err error) {
	readBehaviour := model.ReadBehaviourAgain
	if read.ReadDirection != nil && read.ReadDirection.Previous {
		readBehaviour = model.ReadBehaviourPrevious
	} else if read.ReadDirection != nil && read.ReadDirection.Next {
		readBehaviour = model.ReadBehaviourNext
	}
	if i.Cursor == nil {
		if i.Cursor, err = i.Stream.Storage.InitCursor(i.StartPosition); err != nil {
			return err
		}
	}
	var data []ast.Json
	var endOfStream bool
	i.Cursor, data, endOfStream, err = i.Stream.Storage.Read(readBehaviour, i.Cursor, i.Count)
	for _, d := range data {
		fmt.Println(d)
	}
	if endOfStream {
		fmt.Println("end of the stream")
	}
	return err
}

func NewStreamIteratorFrom(stream *Stream, read *ast.StreamRead) (iterator StreamIterator) {
	startPosition := model.StartPositionBeginning
	if read.ReadClause.StartAtEnd {
		startPosition = model.StartPositionEnd
	}
	count := 10
	if read.ReadClause.Count != nil {
		count = *read.ReadClause.Count
	}
	iterator = StreamIterator{
		Stream:        stream,
		Cursor:        nil,
		StartPosition: startPosition,
		Count:         count,
	}
	return iterator
}
