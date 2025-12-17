package service

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"project_uas/app/repository"
)

type LecturerService struct {
	Repo *repository.LecturerRepo
}

func NewLecturerService(repo *repository.LecturerRepo) *LecturerService {
	return &LecturerService{Repo: repo}
}

func (s *LecturerService) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}
	userIDStr := userID.(string)

	data, err := s.Repo.GetByUserID(userIDStr)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "lecturer not found",
		})
	}

	return c.Status(http.StatusOK).JSON(data)
}

func (s *LecturerService) GetAll(c *fiber.Ctx) error {
	data, err := s.Repo.GetAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed get lecturers",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data": data,
	})
}

