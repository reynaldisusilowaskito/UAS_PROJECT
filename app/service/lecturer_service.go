package service

import (
	"net/http"

	"project_uas/app/repository"
	"github.com/gin-gonic/gin"
)

type LecturerService struct {
	Repo *repository.LecturerRepo
}

func NewLecturerService(repo *repository.LecturerRepo) *LecturerService {
	return &LecturerService{Repo: repo}
}

func (s *LecturerService) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	data, err := s.Repo.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "lecturer not found"})
		return
	}

	c.JSON(http.StatusOK, data)
}
