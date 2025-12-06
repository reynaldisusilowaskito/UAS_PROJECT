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

	// (2) Ambil user berdasarkan username
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

	// (4) Cek status aktif
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "account disabled"})
		return
	}

	// (5) Ambil role
	roleName, err := s.AuthRepo.GetRoleNameByID(user.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot load role"})
		return
	}

	// Ambil permissions user
	perms, err := s.AuthRepo.GetPermissionsByRole(user.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot load permissions"})
		return
	}

	// (6) Generate Access Token
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

	// (7) Return response sesuai SRS
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

// =====================
//     GET PROFILE
// =====================
func (s *AuthService) GetProfile(c *gin.Context) {

	userID := c.GetString("user_id")
	username := c.GetString("username")
	role := c.GetString("role")

	permissions, _ := c.Get("permissions") // c.Get mengembalikan (value, exists)

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user_id":     userID,
			"username":    username,
			"role":        role,
			"permissions": permissions,
		},
	})
}

