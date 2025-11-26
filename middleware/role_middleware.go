package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OnlyAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")

		if role != "Admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func OnlyStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")

		if role != "Mahasiswa" {
			c.JSON(http.StatusForbidden, gin.H{"error": "student only"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func OnlyLecturer() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")

		if role != "Dosen Wali" {
			c.JSON(http.StatusForbidden, gin.H{"error": "lecturer only"})
			c.Abort()
			return
		}

		c.Next()
	}
}
