package config

import (
	"database/sql"
	"time"
)

func OpenDB() *sql.DB {
	var db *sql.DB
	var err error

	db, err = sql.Open("mysql", "root:admin@tcp(localhost:3306)/todo?parseTime=true")
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db
}
