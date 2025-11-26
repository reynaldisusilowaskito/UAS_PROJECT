package service

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"project_uas/app/model"
	"project_uas/app/repository"
)

type AchievementService struct {
	Repo *repository.AchievementRepo
}

func NewAchievementService(repo *repository.AchievementRepo) *AchievementService {
	return &AchievementService{Repo: repo}
}

func (s *AchievementService) CreateAchievement(c *gin.Context) {
	var req model.Achievement
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Wajib punya student_id
	if req.StudentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student_id is required"})
		return
	}

	// 1. Simpan ke MongoDB
	mongoID, err := s.Repo.CreateAchievementMongo(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed saving to mongo"})
		return
	}

	// 2. Simpan ke PostgreSQL 
	ref := model.AchievementReference{
		ID:                 uuid.New().String(),
		StudentID:          req.StudentID,
		MongoAchievementID: mongoID,
		Status:             "draft",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := s.Repo.CreateReference(ref); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed saving reference"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "achievement created",
		"data":    ref,
	})
}
