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

// =====================
// GET /students
// =====================
func (s *StudentService) GetAll(c *fiber.Ctx) error {
	data, err := s.Repo.GetAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}

// =====================
// GET /students/:id
// =====================
func (s *StudentService) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	student, err := s.Repo.GetByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "student not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": student,
	})
}

// =====================
// GET /students/:id/achievements
// =====================
func (s *StudentService) GetAchievements(c *fiber.Ctx) error {
	id := c.Params("id")

	data, err := s.Repo.GetAchievementsByStudentID(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}

// =====================
// PUT /students/:id/advisor (ADMIN)
// =====================
func (s *StudentService) UpdateAdvisor(c *fiber.Ctx) error {
	id := c.Params("id")

	var body struct {
		AdvisorID string `json:"advisor_id"`
	}

	if err := c.BodyParser(&body); err != nil || body.AdvisorID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "advisor_id required",
		})
	}

	if err := s.Repo.UpdateAdvisor(id, body.AdvisorID); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "advisor updated successfully",
	})
}
