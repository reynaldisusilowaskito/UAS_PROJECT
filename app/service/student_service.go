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

func (s *StudentService) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	data, err := s.Repo.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}

	c.JSON(http.StatusOK, data)
}
