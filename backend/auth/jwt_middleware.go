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

// JwtMiddleware authenticates requests done from Mikochi's UI with a JWT
type JwtMiddleware struct {
	jwtSecret              []byte
	invalidatedTokens      map[string]Expirable[struct{}]
	invalidatedTokensMutex sync.RWMutex

}

// Initialize the JwtMiddleware with necessary fields
func NewJwtMiddleware(secret []byte) *JwtMiddleware {
	return &JwtMiddleware{
		jwtSecret:         secret,
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

// isTokenInvalidated checks if a token ID is in the invalidated tokens list
func (j *JwtMiddleware) isTokenInvalidated(jti string) bool {
	j.invalidatedTokensMutex.RLock()
	defer j.invalidatedTokensMutex.RUnlock()

	_, exists := j.invalidatedTokens[jti]
	return exists
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

	claims := Claims{}
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

	if !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
		})
		return
	}

	if claims.ID == "" || j.isTokenInvalidated(claims.ID) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Invalid token",
		})
		return
	}

	if claims.Scope != "mikochi-user" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Invalid token",
		})
		return
	}

	c.Set("jti", claims.ID) // useful for logout
	c.Next()
}

// Cleanup deletes stale data from the invalidatedTokens map avoid unchecked memory growth
func (j *JwtMiddleware) Cleanup() {
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

func parseAuthHeader(header string) (string, error) {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("Invalid header")
	}
	return parts[1], nil
}

