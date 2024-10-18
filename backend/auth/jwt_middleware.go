package auth

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JwtMiddleware struct {
	jwtSecret              []byte
	tokenWhitelist         map[string]string
	tokenWhitelistMutex    sync.Mutex
	invalidatedTokens      map[string]struct{}
	invalidatedTokensMutex sync.RWMutex
}

// Initialize the JwtMiddleware with necessary fields
func NewJwtMiddleware(secret []byte) *JwtMiddleware {
	return &JwtMiddleware{
		jwtSecret:         secret,
		tokenWhitelist:    make(map[string]string),
		invalidatedTokens: make(map[string]struct{}),
	}
}

// AddInvalidatedToken adds a token ID to the invalidated tokens list
func (j *JwtMiddleware) AddInvalidatedToken(jti string) {
	j.invalidatedTokensMutex.Lock()
	defer j.invalidatedTokensMutex.Unlock()
	j.invalidatedTokens[jti] = struct{}{}
	for jti := range j.invalidatedTokens {
		log.Printf("Token ID: %s", jti)
	}
	log.Printf("Token invalidated: %s\n", jti) // Log the invalidated token ID
}

// IsTokenInvalidated checks if a token ID is in the invalidated tokens list
func (j *JwtMiddleware) IsTokenInvalidated(jti string) bool {
	_, exists := j.invalidatedTokens[jti]
	log.Printf("Checking if token is invalidated: %s, exists: %v", jti, exists)
	return exists
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
