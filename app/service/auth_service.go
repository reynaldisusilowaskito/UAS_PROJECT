package service

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"project_uas/helper"
	"project_uas/app/model"
	"project_uas/app/repository"
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
func (s *AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest

	// (1) Validasi input JSON
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid JSON",
		})
	}

	// (2) Ambil user berdasarkan username
	user, err := s.AuthRepo.FindByUsername(req.Username)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid username or password",
		})
	}

	// (3) Validasi password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid username or password",
		})
	}

	// (4) Cek status aktif
	if !user.IsActive {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "account disabled",
		})
	}

	// (5) Ambil role
	roleName, err := s.AuthRepo.GetRoleNameByID(user.RoleID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot load role",
		})
	}

	// Ambil permissions user
	perms, err := s.AuthRepo.GetPermissionsByRole(user.RoleID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot load permissions",
		})
	}

	// (6) Generate Access Token
	accessToken, err := helper.GenerateAccessToken(
		user.ID,
		user.Username,
		roleName,
		perms,
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed generate access token",
		})
	}

	// Generate Refresh Token
	refreshToken, err := helper.GenerateRefreshToken(user.ID, roleName)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed generate refresh token",
		})
	}

	// (7) Return response sesuai SRS
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token":        accessToken,
			"refreshToken": refreshToken,
			"user": fiber.Map{
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
func (s *AuthService) GetProfile(c *fiber.Ctx) error {

	userID := c.Locals("user_id")
	username := c.Locals("username")
	role := c.Locals("role")
	permissions := c.Locals("permissions")

	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"user_id":     userID,
			"username":    username,
			"role":        role,
			"permissions": permissions,
		},
	})
}
