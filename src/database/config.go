package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Connect() *sql.DB {
	log.Println("Initializing SQLite database...")
	db, err := sql.Open("sqlite3", "database.db")

	if err != nil {
		log.Fatal(err)
	}

	table, err := db.Prepare(`CREATE TABLE IF NOT EXISTS videos
	(
		video_id TEXT PRIMARY KEY,
		title TEXT,
		description TEXT,
		upload_initiate_time TEXT,
		upload_status INTEGER,
		upload_end_time TEXT
	)`)


	if err != nil {
		log.Fatal(err)
	}
	defer table.Close()
	table.Exec()

	log.Println("Database initialized.")
	return db
}