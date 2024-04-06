package database

import (
	"celeste/src/model/ast"
	"celeste/src/storage"
)

type Stream struct {
	Name        string
	UpStreams   []*Stream
	DownStreams []*Stream
	Storage     storage.Storage
}

func (s *Stream) Append(data ast.Json) (id string, err error) {
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
