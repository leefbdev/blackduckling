package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func initializeDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./scripts.db")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("DROP TABLE scripts")
	if err != nil {
		log.Printf("Failed to delete scripts table: %v", err)
		return nil, err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS scripts (id INTEGER PRIMARY KEY, output TEXT NOT NULL, status TEXT NOT NULL, repoUrl TEXT NOT NULL)")
	if err != nil {
		return nil, err
	}

	return db, nil
}
