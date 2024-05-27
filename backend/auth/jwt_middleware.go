package auth

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JwtMiddleware struct {
	jwtSecret []byte
	tokenWhitelist map[string]string
	tokenWhitelistMutex sync.Mutex
}

// setWhitelist allows a single-use JWTs (for streams)
// Each token is valid for one route and 24h
func (j *JwtMiddleware) setWhitelist(jti, target string) {
	j.tokenWhitelistMutex.Lock()
	j.tokenWhitelist[jti] = target
	j.tokenWhitelistMutex.Unlock()
}

// CheckAuth is a middleware that will return an error if the request doesn't contain a valid auth token
func (j *JwtMiddleware) CheckAuth(c *gin.Context) {
	encodedToken, err := parseAuthHeader(c.GetHeader("Authorization"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "Invalid Authorization header format",
		})
		return
	}

	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if len(j.jwtSecret) > 0 {
			return j.jwtSecret, nil
		}
		return j.jwtSecret, fmt.Errorf("jwt_secret not set")
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

// CheckStreamAuth is a middleware that will return an error if the request
// doesn't contain a valid single-use auth token passed in the auth query param
func (j *JwtMiddleware) CheckStreamAuth(c *gin.Context) {
	encodedToken := c.Query("auth")

	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(encodedToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if len(j.jwtSecret) > 0 {
			return j.jwtSecret, nil
		}
		return j.jwtSecret, fmt.Errorf("jwt_secret not set")
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"err": "Failed to parse token",
		})
		return
	}

	if !token.Valid || !(j.tokenWhitelist[claims.ID] == c.Param("path")) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid or expired token",
		})
		return
	}

	c.Next()
}
