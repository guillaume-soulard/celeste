package main

import (
	"celeste/src/database"
)

func main() {
	db := database.LoadDatabase()
	var err error
	if err = db.ExecuteCommand(`CREATE STREAM logs STORAGE IN MEMORY`); err != nil {
		panic(err)
	}
	if err = db.ExecuteCommand(`logs < {"field": "value"}`); err != nil {
		panic(err)
	}
	if err = db.ExecuteCommand(`READ START logs`); err != nil {
		panic(err)
	}
}
