package service

import (
	"net/http"

	"project_uas/helper"
	"project_uas/app/model"
	"project_uas/app/repository"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	AuthRepo *repository.AuthRepo
}

func NewAuthService(authRepo *repository.AuthRepo) *AuthService {
	return &AuthService{
		AuthRepo: authRepo,
	}
}

// =====================
//        LOGIN
// =====================
func (s *AuthService) Login(c *gin.Context) {
	var req model.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	// 1. Ambil user berdasarkan username
	user, err := s.AuthRepo.FindByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// 2. Cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// 3. Cek status aktif
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "account is disabled"})
		return
	}

	// 4. Ambil role name (Admin, Mahasiswa, Dosen Wali)
	roleName, err := s.AuthRepo.GetRoleNameByID(user.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed retrieving role"})
		return
	}

	// 5. Ambil permissions berdasarkan role
	perms, err := s.AuthRepo.GetPermissionsByRole(user.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed retrieving permissions"})
		return
	}

	// 6. Generate Access Token (SRS-style)
	accessToken, err := helper.GenerateAccessToken(
		user.ID,
		user.Username,
		roleName,
		perms,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed generating token"})
		return
	}

	// 7. Refresh Token (cukup butuh userID + roleName)
	refreshToken, err := helper.GenerateRefreshToken(user.ID, roleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed generating refresh token"})
		return
	}

	// 8. Return response sesuai SRS
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"token":        accessToken,
			"refreshToken": refreshToken,
			"user": gin.H{
				"id":          user.ID,
				"username":    user.Username,
				"fullName":    user.FullName,
				"role":        roleName,
				"permissions": perms,
			},
		},
	})
}
