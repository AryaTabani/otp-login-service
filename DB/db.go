package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "otp-service.db")
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	createTables()
}

func createTables() {
	createuserstable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		phone_number TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP NOT NULL
	);`

	_, err := DB.Exec(createuserstable)
	if err != nil {
		log.Fatalf("failed to create users table: %v", err)
	}

	createotpstable := `
	CREATE TABLE IF NOT EXISTS otps (
		phone_number TEXT PRIMARY KEY,
		code TEXT NOT NULL,
		expires_at TIMESTAMP NOT NULL
	);`

	_, err = DB.Exec(createotpstable)
	if err != nil {
		log.Fatalf("failed to create otps table: %v", err)
	}
}
