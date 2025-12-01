package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OnlyAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("role") != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func OnlyStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("role") != "student" {
			c.JSON(http.StatusForbidden, gin.H{"error": "students only"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func OnlyLecturer() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("role") != "lecturer" {
			c.JSON(http.StatusForbidden, gin.H{"error": "lecturers only"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ================================
//   SELF ACCESS VALIDATION
// ================================
func OnlySelf() gin.HandlerFunc {
	return func(c *gin.Context) {

		role := c.GetString("role")
		userID := c.GetString("user_id")
		paramID := c.Param("id")

		if role == "admin" {
			c.Next()
			return
		}

		if userID != paramID {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "you can only access your own account",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func OnlyStudentSelf() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.GetString("role") != "student" {
			c.JSON(http.StatusForbidden, gin.H{"error": "students only"})
			c.Abort()
			return
		}

		if c.GetString("user_id") != c.Param("id") {
			c.JSON(http.StatusForbidden, gin.H{"error": "you can only access your own account"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func OnlyLecturerSelf() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.GetString("role") != "lecturer" {
			c.JSON(http.StatusForbidden, gin.H{"error": "lecturers only"})
			c.Abort()
			return
		}

		if c.GetString("user_id") != c.Param("id") {
			c.JSON(http.StatusForbidden, gin.H{"error": "you can only access your own account"})
			c.Abort()
			return
		}

		c.Next()
	}
}
