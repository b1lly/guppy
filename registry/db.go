package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var (
	// DB holds the Guppy database.
	DB *sql.DB
)

func initDB() {
	var err error
	DB, err = sql.Open("mysql", "guppy:gup13@tcp(127.0.0.1:3306)/guppy")
	if err != nil {
		panic(err)
	}
}
