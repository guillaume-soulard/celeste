package storage

import (
	"celeste/src/model"
	"celeste/src/model/ast"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"strconv"
)

var idKey = []byte("_id")

func NewLevelDbStorage(streamName string) (storage Storage, err error) {
	var db *leveldb.DB
	if db, err = leveldb.OpenFile(streamName, nil); err != nil {
		return storage, err
	}
	var idStr []byte
	if idStr, err = db.Get(idKey, nil); err != nil && err != leveldb.ErrNotFound {
		return storage, err
	} else if err == leveldb.ErrNotFound {
		idStr = []byte("0")
	}
	var id int
	if id, err = strconv.Atoi(string(idStr)); err != nil {
		return storage, err
	}
	storage = &LevelDb{
		id: uint64(id),
		db: db,
	}
	return storage, err
}

type LevelDb struct {
	id uint64
	db *leveldb.DB
}

func (f *LevelDb) Append(data ast.Json) (id int64, err error) {
	batch := leveldb.Batch{}
	f.id++
	idToUse := []byte(fmt.Sprintf("%d", f.id))
	batch.Put(idKey, idToUse)
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

func (f *LevelDb) InitCursor(startPosition model.StartPosition) (cursor interface{}, err error) {
	minId := []byte("0")
	maxId := []byte(fmt.Sprintf("%d", f.id))
	it := f.db.NewIterator(&util.Range{
		Start: minId,
		Limit: maxId,
	}, nil)
	if startPosition == model.StartPositionBeginning {
		it.First()
	} else if startPosition == model.StartPositionEnd {
		it.Last()
	}
	cursor = it
	return cursor, err
}

func (f *LevelDb) Read(readBehaviour model.ReadBehaviour, cursor interface{}, count int) (newCursor interface{}, data []ast.Json, endOfStream bool, err error) {
	it := cursor.(iterator.Iterator)
	data = make([]ast.Json, 0, count)
	for {
		value := it.Value()
		var parsed *ast.Ast
		if parsed, err = model.Parse(string(value)); err != nil {
			return newCursor, data, endOfStream, err
		}
		if parsed.Json != nil {
			data = append(data, *parsed.Json)
		}
		if readBehaviour == model.ReadBehaviourNext {
			if !it.Next() {
				endOfStream = true
				break
			}
		} else if readBehaviour == model.ReadBehaviourPrevious {
			if !it.Prev() {
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
