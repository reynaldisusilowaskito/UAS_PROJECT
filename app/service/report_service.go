package service

import (
	"net/http"

	"project_uas/app/repository"
	"github.com/gin-gonic/gin"
)

type ReportService struct {
	Repo *repository.ReportRepo
}

func NewReportService(repo *repository.ReportRepo) *ReportService {
	return &ReportService{Repo: repo}
}

func (s *ReportService) GetAchievementStats(c *gin.Context) {
	data, err := s.Repo.CountAchievementsByStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed report"})
		return
	}

	c.JSON(http.StatusOK, data)
}
