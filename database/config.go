package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// DBConfig holds the database connection details
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// NewDBConfig creates a new DBConfig from environment variables
func NewDBConfig() (*DBConfig, error) {
	return &DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("SSL_MODE"),
	}, nil
}

// Connect establishes a connection to the PostgreSQL database
func Connect(config *DBConfig) (*sql.DB, error) {
	log.Println("Initializing PostgreSQL database...")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Name, config.SSLMode)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established.")

	err = initializeTables(db)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	log.Println("Database initialized.")
	return db, nil
}

func initializeTables(db *sql.DB) error {
	tableQueries := []string{
		`CREATE TABLE IF NOT EXISTS videos (
			video_id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			upload_initiate_time TIMESTAMP,
			upload_status SMALLINT CHECK (upload_status IN (0, 1)),
			upload_end_time TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			user_id TEXT,
			expires_at TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);`,
	}

	for _, query := range tableQueries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to execute query: %s, error: %w", query, err)
		}
	}

	return nil
}
