package config

import (
	"crypto/rand"
	"encoding/base64"
	"os"
)

type FileType struct {
	FileType      string `json:"file_type"`
	FileExtension string `json:"file_extension"`
}
type ConfigResponse struct {
	FileSizeLimit      string     `json:"file_size_limit"`
	SupportedFileTypes []FileType `json:"supported_file_types"`
}

var SecretKey string

func init() {
	SecretKey = os.Getenv("JWT_SECRET_KEY")

	if SecretKey == "" {
		bytes := make([]byte, 32)
		_, _ = rand.Read(bytes)
		SecretKey = base64.StdEncoding.EncodeToString(bytes)
		os.Setenv("JWT_SECRET_KEY", SecretKey)
	}
}
