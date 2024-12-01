package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./game.db")
	if err != nil {
		log.Fatal(err)
	}
	initUser()
	initSession()
}
