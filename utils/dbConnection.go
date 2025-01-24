package utils

import (
	"database/sql"
	"log"
)

func OpenDBConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		log.Fatal(err)
	}

	// Enable foreign key constraints
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Failed to enable foreign key constraints:", err)
	}

	return db
}
