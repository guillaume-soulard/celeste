package storage

import (
	"celeste/src/model"
	"celeste/src/model/ast"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"time"
)

func NewLevelDbStorage(streamName string) (storage Storage, err error) {
	var db *leveldb.DB
	if db, err = leveldb.OpenFile(streamName, nil); err != nil {
		return storage, err
	}
	storage = &LevelDb{
		idGenerator: NewIdGenerator(),
		db:          db,
	}
	return storage, err
}

type LevelDb struct {
	idGenerator IdGenerator
	db          *leveldb.DB
}

func (f *LevelDb) Append(data ast.Json) (id string, err error) {
	batch := leveldb.Batch{}
	id = f.idGenerator.NextId(time.Now())
	idToUse := []byte(id)
	batch.Put(idToUse, []byte(data.ToString()))
	err = f.db.Write(&batch, &opt.WriteOptions{
		Sync: true,
	})
	return id, err
}

func (f *LevelDb) Close() (err error) {
	err = f.db.Close()
	return err
}

type LevelDbCursor struct {
	iterator    iterator.Iterator
	endOfStream bool
}

func (f *LevelDb) InitCursor(startPosition model.StartPosition) (cursor interface{}, err error) {
	minId := []byte(model.MinId)
	maxId := []byte(model.MaxId)
	c := LevelDbCursor{}
	c.iterator = f.db.NewIterator(&util.Range{
		Start: minId,
		Limit: maxId,
	}, nil)
	if startPosition == model.StartPositionBeginning {
		c.endOfStream = !c.iterator.First()
	} else if startPosition == model.StartPositionEnd {
		c.endOfStream = !c.iterator.Last()
	}
	cursor = c
	return cursor, err
}

func (f *LevelDb) Read(readBehaviour model.ReadBehaviour, cursor interface{}, count int) (newCursor interface{}, data []ast.Json, endOfStream bool, err error) {
	it := cursor.(LevelDbCursor)
	data = make([]ast.Json, 0, count)
	if it.endOfStream {
		endOfStream = true
		return newCursor, data, endOfStream, err
	}
	for {
		value := it.iterator.Value()
		var parsed *ast.Ast
		if parsed, err = model.Parse(string(value)); err != nil {
			return newCursor, data, endOfStream, err
		}
		if parsed.Json != nil {
			data = append(data, *parsed.Json)
		}
		if readBehaviour == model.ReadBehaviourNext {
			if !it.iterator.Next() {
				endOfStream = true
				break
			}
		} else if readBehaviour == model.ReadBehaviourPrevious {
			if !it.iterator.Prev() {
				endOfStream = true
				break
			}
		}
		if len(data) >= count {
			break
		}
	}
	return newCursor, data, endOfStream, err
}

func (f *LevelDb) Truncate(evictionPolicies *[]ast.EvictionPolicy) (err error) {
	return err
}
