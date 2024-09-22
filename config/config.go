package config

import (
	"crypto/rand"
	"encoding/base64"
	"os"
)

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
