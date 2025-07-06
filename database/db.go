package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(path string) {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal("failed to open database:", err)
	}

	createTables()
}

func createTables() {
	createPlayerTable := `
	CREATE TABLE IF NOT EXISTS players (
		session_id TEXT PRIMARY KEY,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		score INTEGER DEFAULT 0
	);`

	createFlagsTable := `
	CREATE TABLE IF NOT EXISTS flags (
		session_id TEXT PRIMARY KEY,
		flag TEXT NOT NULL,
		value INTEGER DEFAULT 100
	);`

	createSubmissionsTable := `
	CREATE TABLE IF NOT EXISTS submissions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT,
		flag TEXT,
		submitted_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createServicesTable := `
	CREATE TABLE IF NOT EXISTS services (
		session_id TEXT PRIMARY KEY,
		port INTEGER NOT NULL,
		container_id TEXT NOT NULL
	);`

	_, err := DB.Exec(createPlayerTable)
	if err != nil {
		log.Fatal("failed to create players table:", err)
	}

	_, err = DB.Exec(createFlagsTable)
	if err != nil {
		log.Fatal("failed to create flags table:", err)
	}

	_, err = DB.Exec(createSubmissionsTable)
	if err != nil {
		log.Fatal("failed to create submissions table:", err)
	}

	_, err = DB.Exec(createServicesTable)
	if err != nil {
		log.Fatal("failed to create services table:", err)
	}
}
