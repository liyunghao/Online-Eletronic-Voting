package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var SqliteDB *sql.DB

func Initialize(dbName string) {
	// Try open and connect to database
	var err error
	SqliteDB, err = sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatalf("Open database failed. Something WRONG: %v\n", err)
	}
	log.Println("Open database successfully [" + dbName + "]")
}
