package database

import (
	"celeste/src/model"
	"celeste/src/model/ast"
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

func (db *Database) LookUpStream(name string) (stream *Stream, err error) {
	var exists bool
	if stream, exists = db.Streams[name]; !exists {
		err = errors.New(fmt.Sprintf("stream %s not exists", name))
		return stream, err
	}
	return stream, err
}
