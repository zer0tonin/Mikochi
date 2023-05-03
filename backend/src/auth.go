package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

// checkJWT is a middleware that will return an error if the request doesn't contain a valid auth token
func checkJWT(c *gin.Context) {
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

// generateAuthToken makes a new signed JWT token valid ~1 month
func generateAuthToken(secret []byte) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 730)),
		Issuer: "Mikochi",
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// POST /login
// login takes a username/password pair and returns a JWT if they match the corresponding env vars
func login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := c.BindJSON(&credentials)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "Couldn't deserialize credentials",
		})
		return
	}

	if credentials.Username != username || credentials.Password != password {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"err": "Invalid credentials",
		})
		return
	}

	signedToken, err := generateAuthToken(jwtSecret)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "Failed to generate authentication token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": signedToken,
	})
}

// GET /refresh
// refresh returns a new JWT token (should be called after an auth check)
func refresh(c *gin.Context) {
	signedToken, err := generateAuthToken(jwtSecret)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "Failed to generate authentication token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": signedToken,
	})
}

// TODO: check whitelist against ID claim
// TODO: validate that it is using the right HMAC method

// checkSingleUseJWT is a middleware that will return an error if the request
// doesn't contain a valid single-use auth token passed in the auth query param
func checkSingleUseJWT(c *gin.Context) {
	encodedToken := c.Query("auth")

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

// GET /single-use
// singleUse returns a new single-use JWT token for use in streams
func singleUse(c *gin.Context) {
	jti := uuid.New().String()
	tokenWhitelist[jti] = true

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		Issuer: "Mikochi",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ID: jti,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "Failed to generate authentication token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": signedToken,
	})
}
