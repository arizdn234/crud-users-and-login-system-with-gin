package middleware

import (
	"net/http"

	"github.com/arizdn234/crud-users-and-login-system-with-gin/utils"
	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken, err := c.Cookie("auth_token")
		if err != nil || authToken == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenClaims, err := utils.VerifyToken(authToken)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("id", tokenClaims.ID)

		c.Next()
	}
}
