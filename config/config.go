package config

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
)

type FileType struct {
	FileType      string `json:"file_type"`
	FileExtension string `json:"file_extension"`
}
type ConfigResponse struct {
	FileSizeLimit      string     `json:"file_size_limit"`
	SupportedFileTypes []FileType `json:"supported_file_types"`
}

type Config struct {
	RootPath               string
	AppwriteBucketID       string
	AppwriteProjectID      string
	AppwriteKey            string
	AppwriteResponseFormat string
	DBHost                 string
	DBPort                 string
	DBUser                 string
	DBPassword             string
	DBName                 string
	Port                   string
	Addr                   string
	SSLMode                string
	JWTSecretKey           string
	FileSizeLimit          string
}

var AppConfig *Config

func LoadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error reading %s: %w", filename, err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			return fmt.Errorf("error reading invalid environment variable definition: %s", line)
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		os.Setenv(key, value)
		log.Printf("Set %s=%s", key, value)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading %s: %w", filename, err)
	}

	return nil
}

func LoadConfig(envFile string) (*Config, error) {
	if err := LoadEnvFile(envFile); err != nil {
		return nil, err
	}

	config := &Config{
		RootPath:               os.Getenv("ROOT_PATH"),
		AppwriteBucketID:       os.Getenv("BUCKET_ID"),
		AppwriteProjectID:      os.Getenv("APPWRITE_PROJECT_ID"),
		AppwriteKey:            os.Getenv("APPWRITE_KEY"),
		AppwriteResponseFormat: os.Getenv("APPWRITE_RESPONSE_FORMAT"),
		DBHost:                 os.Getenv("DB_HOST"),
		DBPort:                 os.Getenv("DB_PORT"),
		DBUser:                 os.Getenv("DB_USER"),
		DBPassword:             os.Getenv("DB_PASSWORD"),
		DBName:                 os.Getenv("DB_NAME"),
		Port:                   os.Getenv("PORT"),
		Addr:                   os.Getenv("ADDR"),
		SSLMode:                os.Getenv("SSL_MODE"),
		JWTSecretKey:           os.Getenv("JWT_SECRET_KEY"),
		FileSizeLimit:          os.Getenv("FILE_SIZE_LIMIT"),
	}

	if config.JWTSecretKey == "" {
		bytes := make([]byte, 32)
		_, _ = rand.Read(bytes)
		os.Setenv("JWT_SECRET_KEY", base64.StdEncoding.EncodeToString(bytes))
		config.JWTSecretKey = os.Getenv("FILE_SIZE_LIMIT")
	}

	AppConfig = config

	log.Println("Configuration loaded successfully.")
	return config, nil
}
