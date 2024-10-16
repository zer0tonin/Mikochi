package auth

import (
    "log"
    "strings"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)

type AuthHandlers struct {
    authMiddleware AuthMiddleware
    username       string
    password       string
    jwtSecret      []byte
    jwtMiddleware  *JwtMiddleware // Add jwtMiddleware field
}

func NewAuthHandlers(authMiddleware AuthMiddleware, username, password string, jwtSecret []byte, jwtMiddleware *JwtMiddleware) *AuthHandlers {
    return &AuthHandlers{
        authMiddleware: authMiddleware,
        username:       username,
        password:       password,
        jwtSecret:      jwtSecret,
        jwtMiddleware:  jwtMiddleware, // Initialize jwtMiddleware
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
    signedToken, err := generateAuthToken(a.jwtSecret)
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
    signedToken, err := generateAuthToken(a.jwtSecret)
    
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
    

    a.authMiddleware.setWhitelist(jti, c.Query("target"))

    claims := jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
        Issuer:    "Mikochi",
        IssuedAt:  jwt.NewNumericDate(time.Now()),
        ID:        jti,
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(a.jwtSecret)

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

// POST /logout
// Logout invalidates the current JWT token
func (h *AuthHandlers) Logout(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid authorization header"})
        return
    }

    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    log.Printf("Token: %s\n", tokenString)

    if tokenString == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
        return
    }

    claims := &jwt.RegisteredClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return h.jwtSecret, nil
    })
    if err != nil || !token.Valid {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
        return
    }

    // Log the entire claims object for debugging
    log.Printf("Extracted claims: %+v\n", claims)

    log.Printf("Extracted jti: %s\n", claims.ID) // Log the extracted jti
    h.jwtMiddleware.AddInvalidatedToken(claims.ID)
    c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully", "invalidated_token": claims.ID})
}