package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize SQLite database
	db, err := sql.Open("sqlite3", "./data/blog.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Create tables
	query := `
	CREATE TABLE IF NOT EXISTS authors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		password TEXT NOT NULL,
		is_admin BOOLEAN DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		slug TEXT NOT NULL UNIQUE,
		body TEXT NOT NULL,
		metadata TEXT NOT NULL,
		deleted BOOLEAN DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		author_id INTEGER REFERENCES authors(id) NOT NULL
	);
	`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	fmt.Println("Database initialized successfully!")
}