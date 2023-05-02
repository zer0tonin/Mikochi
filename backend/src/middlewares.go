package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func parseAuthHeader(header string) (string, error) {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("Invalid header")
	}
	return parts[1], nil
}

func checkJWT(c *gin.Context) {
	encodedToken, err := parseAuthHeader(c.GetHeader("Authorization"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "Invalid Authorization header format",
		})
		return
	}

	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		secret := []byte(os.Getenv("JWT_SECRET"))
		if len(secret) > 0 {
			return secret, nil
		}
		return secret, fmt.Errorf("JWT_SECRET not set")
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"err": "Failed to parse token",
		})
		return
	}

	if !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid or expired token",
		})
		return
	}
	c.Next()
}
