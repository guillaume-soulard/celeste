package database

import (
	"celeste/src/model"
	"celeste/src/model/ast"
	"celeste/src/storage"
	"errors"
	"fmt"
)

func LoadDatabase() Database {
	return Database{
		Streams: make(map[string]*Stream),
	}
}

type Database struct {
	Streams map[string]*Stream
}

func (db *Database) ExecuteCommand(command string) (err error) {
	var parsed *ast.Ast
	if parsed, err = model.Parse(command); err != nil {
		return err
	}
	found := false
	for _, executor := range CommandExecutors {
		if executor.Accept(parsed) {
			found = true
			if err = executor.Execute(db, parsed); err != nil {
				return err
			}
		}
	}

	if !found {
		err = errors.New(fmt.Sprintf("no implementation for command %s ", command))
	}
	return err
}

func (db *Database) NewStreamFrom(creation ast.StreamCreation) (err error) {
	if _, streamAlreadyExists := db.Streams[*creation.Name]; streamAlreadyExists {
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
	if creation.StreamDataSource != nil && creation.StreamDataSource.From != nil {
		if upStream, exists := db.Streams[*creation.StreamDataSource.From]; exists {
			upStream.DownStreams = append(upStream.DownStreams, &stream)
			stream.UpStreams = append(stream.UpStreams, upStream)
		} else {
			err = errors.New(fmt.Sprintf("datasource stream %s not exists", *creation.StreamDataSource.From))
			return err
		}
	}
	db.Streams[stream.Name] = &stream
	return err
}

func (db *Database) LookUpStream(name string) (stream *Stream, err error) {
	var exists bool
	if stream, exists = db.Streams[name]; !exists {
		err = errors.New(fmt.Sprintf("stream %s not exists", name))
		return stream, err
	}
	return stream, err
}

func (db *Database) DropStream(streamName string) (err error) {
	if _, err = db.LookUpStream(streamName); err != nil {
		return err
	}
	delete(db.Streams, streamName)
	return err
}
