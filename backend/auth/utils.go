package auth

import (
	"crypto/rand"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func parseAuthHeader(header string) (string, error) {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("Invalid header")
	}
	return parts[1], nil
}

// generateAuthToken makes a new signed JWT token valid ~1 month
func generateAuthToken(secret []byte) (string, error) {
	jti := uuid.New().String() // Generate a unique jti
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 730)),
		Issuer:    "Mikochi",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        jti, // Add the jti claim
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

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
