package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JwtMiddleware struct {
	jwtSecret              []byte

	tokenWhitelist         map[string]Expirable[string]
	tokenWhitelistMutex    sync.RWMutex

	invalidatedTokens      map[string]Expirable[struct{}]
	invalidatedTokensMutex sync.RWMutex
}

// Initialize the JwtMiddleware with necessary fields
func NewJwtMiddleware(secret []byte) *JwtMiddleware {
	return &JwtMiddleware{
		jwtSecret:         secret,
		tokenWhitelist:    make(map[string]Expirable[string]),
		invalidatedTokens: make(map[string]Expirable[struct{}]),
	}
}

// AddInvalidatedToken adds a token ID to the invalidated tokens list
func (j *JwtMiddleware) InvalidateToken(c *gin.Context) error {
	jti, exists := c.Get("jti")
	if !exists {
		return fmt.Errorf("No JIT in context")
	}
	jtiStr, ok := jti.(string)
	if !ok {
		return fmt.Errorf("JIT is not a string")
	}

	j.invalidatedTokensMutex.Lock()
	defer j.invalidatedTokensMutex.Unlock()

	// struct{}{} = 0 size value
	j.invalidatedTokens[jtiStr] = NewExpirable(struct{}{}, time.Hour * 730)
	log.Printf("Token invalidated: %s\n", jti)
	return nil
}

// IsTokenInvalidated checks if a token ID is in the invalidated tokens list
func (j *JwtMiddleware) IsTokenInvalidated(jti string) bool {
	j.invalidatedTokensMutex.RLock()
	defer j.invalidatedTokensMutex.RUnlock()

	_, exists := j.invalidatedTokens[jti]
	return exists
}

// setWhitelist assigns a JWT ID to the file path it is valid for
// Each token is valid for only one file and 24h
// This is used for streams
func (j *JwtMiddleware) setWhitelist(jti, target string) {
	j.tokenWhitelistMutex.Lock()
	defer j.tokenWhitelistMutex.Unlock()

	j.tokenWhitelist[jti] = NewExpirable(target, 24 * time.Hour)
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

	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (any, error) {
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

	// Check if the token is invalidated and set the jti in the context
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if jti, ok := claims["jti"].(string); ok {
			if j.IsTokenInvalidated(jti) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "Token has been invalidated",
				})
				return
			}
			c.Set("jti", jti)
		}
	}

	c.Next()
}

// CheckStreamAuth is a middleware that will return an error if the request
// doesn't contain a valid stream auth token matching the file being requested
// passed in the auth query param
func (j *JwtMiddleware) CheckStreamAuth(c *gin.Context) {
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

func (j *JwtMiddleware) cleanupStreamWhitelist() {
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

func (j *JwtMiddleware) cleanupInvalidatedTokens() {
	j.invalidatedTokensMutex.Lock()
	defer j.invalidatedTokensMutex.Unlock()

	expired := []string{}
	for key, value := range j.invalidatedTokens {
		if value.IsExpired() {
			expired = append(expired, key)
		}
	}

	for _, e := range expired {
		delete(j.invalidatedTokens, e)
	}
}

// Cleanup deletes stale data from the tokenWhitelist to avoid unchecked memory growth
func (j *JwtMiddleware) Cleanup() {
	j.cleanupStreamWhitelist()
}

func parseAuthHeader(header string) (string, error) {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("Invalid header")
	}
	return parts[1], nil
}

