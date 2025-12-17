package service

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"project_uas/app/repository"
)

type ReportService struct {
	Repo *repository.ReportRepo
}

func NewReportService(repo *repository.ReportRepo) *ReportService {
	return &ReportService{Repo: repo}
}

func (s *ReportService) GetAchievementStats(c *fiber.Ctx) error {
	data, err := s.Repo.CountAchievementsByStatus()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed report",
		})
	}

	return c.Status(http.StatusOK).JSON(data)
}

// GET /api/v1/reports/student/:id
func (s *ReportService) GetStudentReport(c *fiber.Ctx) error {
	studentID := c.Params("id")

	data, err := s.Repo.GetStudentAchievementReport(studentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed student report",
		})
	}

	return c.JSON(fiber.Map{
		"student_id": studentID,
		"data":       data,
	})
}
