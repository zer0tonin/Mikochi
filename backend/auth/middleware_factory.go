package auth

import (
	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	CheckAuth(c *gin.Context)
	InvalidateToken(c *gin.Context) error
	Cleanup()
}

func NewAuthMiddleware(withAuth bool, jwtSecret string) AuthMiddleware {
	if withAuth {
		return NewJwtMiddleware([]byte(jwtSecret))
	}
	return &NoauthMiddleware{}
}

type StreamAuthMiddleware interface {
	CheckAuth(c *gin.Context)
}

func NewStreamAuthMiddleware(withAuth bool, jwtSecret string) StreamAuthMiddleware {
	if withAuth {
		return NewStreamJwtMiddleware([]byte(jwtSecret))
	}
	return &NoauthMiddleware{}
}
