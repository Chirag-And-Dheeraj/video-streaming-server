package database

import (
	"database/sql"
	"fmt"
	"log"
	"video-streaming-server/config"

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
		Host:     config.AppConfig.DBHost,
		Port:     config.AppConfig.DBPort,
		User:     config.AppConfig.DBUser,
		Password: config.AppConfig.DBPassword,
		Name:     config.AppConfig.DBName,
		SSLMode:  config.AppConfig.SSLMode,
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
