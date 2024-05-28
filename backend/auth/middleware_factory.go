package auth

import (
	"sync"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	setWhitelist(jti, target string)
	CheckAuth(c *gin.Context)
	CheckStreamAuth(c *gin.Context)
}

func NewAuthMiddleware(withAuth bool, jwtSecret string) AuthMiddleware {
	if withAuth {
		return &JwtMiddleware{
			jwtSecret:           []byte(jwtSecret),
			tokenWhitelist:      map[string]string{},
			tokenWhitelistMutex: sync.Mutex{},
		}
	}
	return &NoauthMiddleware{}
}
