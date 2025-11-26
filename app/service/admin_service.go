package service

import (
	"net/http"
	"project_uas/app/repository"
	"github.com/gin-gonic/gin"
)

type AdminService struct {
	Repo *repository.AdminRepo
}

func NewAdminService(repo *repository.AdminRepo) *AdminService {
	return &AdminService{Repo: repo}
}

func (s *AdminService) GetAllUsers(c *gin.Context) {
	data, err := s.Repo.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed get users"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (s *AdminService) UpdateUserStatus(c *gin.Context) {
	type Req struct {
		UserID string `json:"user_id"`
		Active bool   `json:"active"`
	}

	var req Req
	c.ShouldBindJSON(&req)

	err := s.Repo.UpdateUserStatus(req.UserID, req.Active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}
