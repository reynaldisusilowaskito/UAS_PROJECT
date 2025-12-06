	package middleware

	import (
		"net/http"
		"strings"

		"project_uas/helper"

		"github.com/gin-gonic/gin"
	)

	func AuthMiddleware() gin.HandlerFunc {
		return func(c *gin.Context) {

			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
				c.Abort()
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header"})
				c.Abort()
				return
			}

			accessToken := parts[1]

			claims, err := helper.VerifyAccessToken(accessToken)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
				c.Abort()
				return
			}

			// Set ke context
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("role", claims.Role)
			c.Set("permissions", claims.Permissions)

			c.Next()
		}
	}
