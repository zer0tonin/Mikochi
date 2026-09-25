package auth

import (
	"crypto/rand"
	"log"
)

// GenerateRandomSecret creates a 256 bytes array used as a jwt secret when no env var is set
func GenerateRandomSecret() string {
	bytes := make([]byte, 256)
	_, err := rand.Read(bytes)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		log.Panicf("Failed to generate JWT secret: %s", err.Error())
	}
	return string(bytes)
}
