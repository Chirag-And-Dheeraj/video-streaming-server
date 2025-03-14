package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

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
func newDBConfig() (*DBConfig, error) {
	return &DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("SSL_MODE"),
	}, nil
}

func initializeTables(db *sql.DB) error {
	tableQueries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS videos (
			video_id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			upload_initiate_time TIMESTAMP,
			upload_status SMALLINT CHECK (upload_status IN (0, 1)),
			upload_end_time TIMESTAMP,
			user_id TEXT,
			delete_flag SMALLINT CHECK (delete_flag IN (0, 1)),
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`,
		`ALTER TABLE videos
			ADD COLUMN IF NOT EXISTS thumbnail TEXT;`,
	}

	for _, query := range tableQueries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to execute query: %s, error: %w", query, err)
		}
	}

	return nil
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


	log.Println("Database initialized.")
	return db, nil
}

func GetDBConn() *sql.DB {
	if DB == nil {
		dbConfig, err := newDBConfig()
		if err != nil {
			log.Fatalf("Failed to load database config: %v", err)
		}
		DB, err = Connect(dbConfig)

		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
	}
	return DB
}
