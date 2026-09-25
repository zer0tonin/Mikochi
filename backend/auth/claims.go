package auth

import "github.com/golang-jwt/jwt/v5"

// Custom JWT claims
type Claims struct {
	Scope string `json:"scope,omitempty"`
	jwt.RegisteredClaims
}
