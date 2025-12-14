package service

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"project_uas/app/repository"
)

type StudentService struct {
	Repo *repository.StudentRepo
}

func NewStudentService(repo *repository.StudentRepo) *StudentService {
	return &StudentService{Repo: repo}
}

// GET /students/profile
func (s *StudentService) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing user in token",
		})
	}
	userIDStr := userID.(string)

	student, err := s.Repo.FindByUserID(userIDStr)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "student profile not found",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "student profile loaded",
		"data":    student,
	})
}
