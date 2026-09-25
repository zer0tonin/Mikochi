package auth

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// StreamJwtMiddleware authenticates requests done when streaming files (from VLC, etc...)
type StreamJwtMiddleware struct {
	jwtSecret              []byte
	tokenWhitelist         map[string]Expirable[string]
	tokenWhitelistMutex    sync.RWMutex
}

// Initialize the JwtMiddleware with necessary fields
func NewStreamJwtMiddleware(secret []byte) *StreamJwtMiddleware {
	return &StreamJwtMiddleware{
		jwtSecret:         secret,
		tokenWhitelist:    make(map[string]Expirable[string]),
	}
}

// setWhitelist assigns a JWT ID to the file path it is valid for
// Each token is valid for only one file and 24h
// This is used for streams
func (j *StreamJwtMiddleware) setWhitelist(jti, target string) {
	j.tokenWhitelistMutex.Lock()
	defer j.tokenWhitelistMutex.Unlock()

	j.tokenWhitelist[jti] = NewExpirable(target, 24 * time.Hour)
}

// CheckStreamAuth is a middleware that will return an error if the request
// doesn't contain a valid stream auth token matching the file being requested
// passed in the auth query param
func (j *StreamJwtMiddleware) CheckAuth(c *gin.Context) {
	encodedToken := c.Query("auth")

	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(encodedToken, &claims, func(token *jwt.Token) (any, error) {
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

	if !token.Valid || claims.ID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
		})
		return
	}

	j.tokenWhitelistMutex.RLock()
	expirable, ok := j.tokenWhitelist[claims.ID]
	j.tokenWhitelistMutex.RUnlock()

	if !ok || expirable.IsExpired() || expirable.GetValue() != c.Param("path") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Expired token",
		})
		return
	}

	c.Next()
}

func (j *StreamJwtMiddleware) cleanupStreamWhitelist() {
	j.tokenWhitelistMutex.Lock()
	defer j.tokenWhitelistMutex.Unlock()

	expired := []string{}
	for key, value := range j.tokenWhitelist {
		if value.IsExpired() {
			expired = append(expired, key)
		}
	}

	for _, e := range expired {
		delete(j.tokenWhitelist, e)
	}
}

// Cleanup deletes stale data from the tokenWhitelist to avoid unchecked memory growth
func (j *StreamJwtMiddleware) Cleanup() {
	j.cleanupStreamWhitelist()
}
