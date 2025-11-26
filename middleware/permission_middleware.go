package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {

		perms, exists := c.Get("permissions")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "permissions not found in token"})
			c.Abort()
			return
		}

		permList, ok := perms.([]string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid permission format"})
			c.Abort()
			return
		}

		allowed := false
		for _, p := range permList {
			if p == permission {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "you do not have permission: " + permission,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
