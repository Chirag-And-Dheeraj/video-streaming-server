package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func Connect() *sql.DB {
	log.Println("Initializing PostgreSQL database...")

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("SSL_MODE")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}

	query, err := db.Prepare(`CREATE TABLE IF NOT EXISTS videos (
		video_id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		upload_initiate_time TIMESTAMP,
		upload_status SMALLINT CHECK (upload_status IN (0, 1)),
		upload_end_time TIMESTAMP
	);`)

	if err != nil {
		log.Fatal(err)
	}

	defer query.Close()
	query.Exec()

	log.Println("Database initialized.")
	return db
}
