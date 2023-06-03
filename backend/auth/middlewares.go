package auth

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// Whitelist of single-use JWTs (for streams)
// Each token is valid for one route and 24h
var tokenWhitelist = map[string]string{}
var tokenWhitelistMutex = sync.Mutex{}

// CheckJWT is a middleware that will return an error if the request doesn't contain a valid auth token
func CheckJWT(c *gin.Context) {
	jwtSecret := []byte(viper.GetString("JWT_SECRET"))
	encodedToken, err := parseAuthHeader(c.GetHeader("Authorization"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "Invalid Authorization header format",
		})
		return
	}

	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if len(jwtSecret) > 0 {
			return jwtSecret, nil
		}
		return jwtSecret, fmt.Errorf("jwt_secret not set")
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

// CheckSingleUseJWT is a middleware that will return an error if the request
// doesn't contain a valid single-use auth token passed in the auth query param
func CheckSingleUseJWT(c *gin.Context) {
	jwtSecret := []byte(viper.GetString("JWT_SECRET"))
	encodedToken := c.Query("auth")

	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(encodedToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if len(jwtSecret) > 0 {
			return jwtSecret, nil
		}
		return jwtSecret, fmt.Errorf("jwt_secret not set")
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"err": "Failed to parse token",
		})
		return
	}

	if !token.Valid || !(tokenWhitelist[claims.ID] == c.Param("path")) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid or expired token",
		})
		return
	}

	c.Next()
}
