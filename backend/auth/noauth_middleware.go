package auth

import "github.com/gin-gonic/gin"

type NoauthMiddleware struct {
}

func (n *NoauthMiddleware) CheckAuth(c *gin.Context) {
	c.Next()
}

func (n *NoauthMiddleware) CheckStreamAuth(c *gin.Context) {
	c.Next()
}

func (n *NoauthMiddleware) InvalidateToken(c *gin.Context) error {
	return nil
}

func (n *NoauthMiddleware) Cleanup() {}
