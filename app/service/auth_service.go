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

	// (1) Validasi input JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	// (2) Ambil user berdasarkan username/email
	user, err := s.AuthRepo.FindByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// (3) Validasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// Cek status user
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "account disabled"})
		return
	}

	// (4) Ambil role user
	roleName, err := s.AuthRepo.GetRoleNameByID(user.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get user role"})
		return
	}

	// Ambil permissions user
	perms, err := s.AuthRepo.GetPermissionsByRole(user.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get permissions"})
		return
	}

	// (5) Generate Access Token
	accessToken, err := helper.GenerateAccessToken(
		user.ID,
		user.Username,
		roleName,
		perms,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed generate access token"})
		return
	}

	// Generate Refresh Token
	refreshToken, err := helper.GenerateRefreshToken(user.ID, roleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed generate refresh token"})
		return
	}

	// (6) Return response sesuai SRS
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"token":        accessToken,
			"refreshToken": refreshToken,
			"user": gin.H{
				"id":          user.ID,
				"username":    user.Username,
				"email":       user.Email,
				"fullName":    user.FullName,
				"role":        roleName,
				"permissions": perms,
			},
		},
	})
}
