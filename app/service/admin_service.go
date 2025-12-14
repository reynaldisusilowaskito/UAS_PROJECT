package service

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"project_uas/app/repository"
)

type AdminService struct {
	Repo *repository.AdminRepo
}

func NewAdminService(repo *repository.AdminRepo) *AdminService {
	return &AdminService{Repo: repo}
}

func (s *AdminService) GetAllUsers(c *fiber.Ctx) error {
	data, err := s.Repo.GetAllUsers()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed get users",
		})
	}

	return c.Status(http.StatusOK).JSON(data)
}

func (s *AdminService) UpdateUserStatus(c *fiber.Ctx) error {
	type Req struct {
		UserID string `json:"user_id"`
		Active bool   `json:"active"`
	}

	var req Req
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	err := s.Repo.UpdateUserStatus(req.UserID, req.Active)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed update user",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "updated",
	})
}
