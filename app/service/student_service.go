package service

import (
	"net/http"

	"project_uas/app/repository"
	"github.com/gin-gonic/gin"
)

type StudentService struct {
	Repo *repository.StudentRepo
}

func NewStudentService(repo *repository.StudentRepo) *StudentService {
	return &StudentService{Repo: repo}
}

// GET /students/profile
func (s *StudentService) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user in token"})
		return
	}

	student, err := s.Repo.FindByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "student profile not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "student profile loaded",
		"data":    student,
	})
}
