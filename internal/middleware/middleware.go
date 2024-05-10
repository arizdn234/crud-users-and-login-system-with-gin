package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("auth_token")
		if err != nil || cookie == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
