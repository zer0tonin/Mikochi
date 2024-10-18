package auth

import (
	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	setWhitelist(jti, target string)
	CheckAuth(c *gin.Context)
	CheckStreamAuth(c *gin.Context)
	InvalidateToken(c *gin.Context) error
}

func NewAuthMiddleware(withAuth bool, jwtSecret string) AuthMiddleware {
	if withAuth {
		return NewJwtMiddleware([]byte(jwtSecret))
	}
	return &NoauthMiddleware{}
}
