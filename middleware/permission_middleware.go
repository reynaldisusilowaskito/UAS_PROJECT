package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {

		perms := c.GetStringSlice("permissions")

		for _, p := range perms {
			if p == permission {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "missing permission: " + permission,
		})
		c.Abort()
	}
}
