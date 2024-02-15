package database

import (
	"celeste/src/model/ast"
	"celeste/src/storage"
	"errors"
	"fmt"
)

type Stream struct {
	Name        string
	UpStreams   []*Stream
	DownStreams []*Stream
	Storage     storage.Storage
}

func NewStreamFrom(database *Database, creation ast.StreamCreation) (err error) {
	if _, streamAlreadyExists := database.Streams[*creation.Name]; streamAlreadyExists {
		err = errors.New(fmt.Sprintf("stream %s already exists", *creation.Name))
		return err
	}
	var s storage.Storage
	if s, err = storage.NewStorageFrom(creation); err != nil {
		return err
	}
	stream := Stream{
		Name:        *creation.Name,
		UpStreams:   make([]*Stream, 0, 10),
		DownStreams: make([]*Stream, 0, 10),
		Storage:     s,
	}
	database.Streams[stream.Name] = &stream
	return err
}

func (s *Stream) Append(data ast.Json) (id int64, err error) {
	if id, err = s.Storage.Append(data); err != nil {
		return id, err
	}
	for _, stream := range s.DownStreams {
		if _, err = stream.Append(data); err != nil {
			return id, err
		}
	}
	return id, err
}

func (s *Stream) Iterator(read *ast.StreamRead) (iterator StreamIterator, err error) {
	iterator = NewStreamIteratorFrom(s, read)
	return iterator, err
}
