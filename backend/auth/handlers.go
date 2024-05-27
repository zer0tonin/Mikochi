package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type AuthWhitelist interface {
	setWhitelist(jti, target string)
}

type AuthHandlers struct {
	authWhitelist AuthWhitelist
	username string
	password string
}

func NewAuthHandlers(authWhitelist AuthWhitelist, username, password string) *AuthHandlers {
	return &AuthHandlers{
		authWhitelist: authWhitelist,
	}
}

// POST /login
// Login takes a username/password pair and returns a JWT if they match the corresponding env vars
func (a *AuthHandlers) Login(c *gin.Context) {
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

	ok := rateLimiter.checkRateLimit(credentials.Username)
	if !ok {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"err": "You are being rate limited",
		})
		return
	}

	if credentials.Username != a.username || credentials.Password != a.password {
		log.Printf("Failed login attempt for %s", credentials.Username)
		rateLimiter.increaseRateLimit(credentials.Username)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"err": "Invalid credentials",
		})
		return
	}

	log.Printf("Succesful login as %s", credentials.Username)
	rateLimiter.resetRateLimit(credentials.Username)
	signedToken, err := generateAuthToken([]byte(viper.GetString("JWT_SECRET")))
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
// Refresh returns a new JWT token (should be called after an auth check)
func (a *AuthHandlers) Refresh(c *gin.Context) {
	signedToken, err := generateAuthToken([]byte(viper.GetString("JWT_SECRET")))
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

// GET /single-use
// SingleUse returns a new single-use JWT token for use in streams
func (a *AuthHandlers) SingleUse(c *gin.Context) {
	jti := uuid.New().String()


	a.authWhitelist.setWhitelist(jti, c.Query("target"))

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		Issuer:    "Mikochi",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        jti,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(viper.GetString("JWT_SECRET")))

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
