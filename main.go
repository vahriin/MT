package main

import "github.com/vahriin/MT/db"

func main() {
	appdb, err := db.InitDB("dbname=MT_DB sslmode=disable")
	if err != nil {
		panic(err)
	}

	appdb.CreateTables()
}
