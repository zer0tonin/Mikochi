package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// StreamJwtMiddleware authenticates requests done when streaming files (from VLC, etc...)
type StreamJwtMiddleware struct {
	jwtSecret              []byte
}

// Initialize the JwtMiddleware with necessary fields
func NewStreamJwtMiddleware(secret []byte) *StreamJwtMiddleware {
	return &StreamJwtMiddleware{
		jwtSecret:         secret,
	}
}

// CheckStreamAuth is a middleware that will return an error if the request
// doesn't contain a valid stream auth token matching the file being requested
// passed in the auth query param
func (j *StreamJwtMiddleware) CheckAuth(c *gin.Context) {
	encodedToken := c.Query("auth")

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

	if !token.Valid || claims.ID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
		})
		return
	}

	if claims.Scope != c.Param("path") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Expired token",
		})
		return
	}

	c.Next()
}
