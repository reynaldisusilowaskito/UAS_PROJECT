package service

import (
	"net/http"
	"log"
	"time"
	"strconv"
	"fmt"
	"path/filepath"
	"strings"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"project_uas/app/model"
	"project_uas/app/repository"
)

type AchievementService struct {
	Repo        *repository.AchievementRepo
	StudentRepo *repository.StudentRepo
}

func NewAchievementService(repo *repository.AchievementRepo, studentRepo *repository.StudentRepo) *AchievementService {
	return &AchievementService{
		Repo:        repo,
		StudentRepo: studentRepo,
	}
}

// ------------------------- CREATE -------------------------
func (s *AchievementService) CreateAchievement(c *fiber.Ctx) error {
	var req model.Achievement

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing user in token",
		})
	}
	userIDStr := userID.(string)

	log.Println("DEBUG Token user_id =", userIDStr)

	student, err := s.StudentRepo.FindByUserID(userIDStr)
	if err != nil {
		log.Println("DEBUG StudentRepo.FindByUserID error =", err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":  "student profile not found",
			"detail": err.Error(),
		})
	}

	log.Println("DEBUG FOUND STUDENT:", student.ID, student.UserID)

	req.CreatedBy = userIDStr
	now := time.Now()
	req.CreatedAt = now
	req.UpdatedAt = now

	s.Repo.EnsureDBs()

	mongoHex, err := s.Repo.CreateAchievementMongo(req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":  "failed insert achievement into mongo",
			"detail": err.Error(),
		})
	}

	ref := model.AchievementReference{
		ID:                 uuid.New().String(),
		StudentID:          student.ID,
		MongoAchievementID: mongoHex,
		Status:             "draft",
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.Repo.CreateReference(ref); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":  "failed insert reference to postgres",
			"detail": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "achievement created",
		"data":    ref,
	})
}

// ------------------------- DETAIL ----------------------------
func (s *AchievementService) GetAchievementDetail(c *fiber.Ctx) error {
	refID := c.Params("id")
	if refID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "missing id"})
	}

	s.Repo.EnsureDBs()

	ref, err := s.Repo.GetReferenceByID(refID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "reference not found"})
	}

	ach, err := s.Repo.GetAchievementMongo(ref.MongoAchievementID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":  "failed retrieve achievement",
			"detail": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"reference":   ref,
		"achievement": ach,
	})
}

// ----------------------- SUBMIT ----------------------------
func (s *AchievementService) SubmitAchievement(c *fiber.Ctx) error {
	refID := c.Params("id")
	if refID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "missing id"})
	}

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "missing user in token"})
	}
	userIDStr := userID.(string)

	s.Repo.EnsureDBs()

	student, err := s.StudentRepo.FindByUserID(userIDStr)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "student profile not found"})
	}

	ref, err := s.Repo.GetReferenceByID(refID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "reference not found"})
	}

	if ref.StudentID != student.ID {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "not allowed"})
	}

	if ref.Status != "draft" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "only draft can be submitted"})
	}

	if err := s.Repo.Submit(refID); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":  "failed to submit",
			"detail": err.Error(),
		})
	}

	_ = s.Repo.AddHistory(refID, "draft", "submitted", userIDStr, "")

	if student.AdvisorID != nil {
		_ = s.Repo.CreateNotification(
			*student.AdvisorID,
			"Pengajuan Prestasi Baru",
			"Mahasiswa mengirim pengajuan prestasi untuk diverifikasi.",
		)
	}

	return c.JSON(fiber.Map{
		"message": "submitted",
		"status":  "submitted",
	})
}

// ------------------------- VERIFY ----------------------------
func (s *AchievementService) VerifyAchievement(c *fiber.Ctx) error {
	refID := c.Params("id")
	if refID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "missing id"})
	}

	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	if role != "lecturer" && role != "admin" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "only lecturer or admin can verify"})
	}

	ref, err := s.Repo.GetReferenceByID(refID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "reference not found"})
	}

	if ref.Status != "submitted" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "only submitted achievements can be verified"})
	}

	student, err := s.StudentRepo.GetByID(ref.StudentID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "student not found"})
	}

	if role == "lecturer" {
		lecturerID, err := s.StudentRepo.GetLecturerIDByUserID(userID)
		if err != nil {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "lecturer profile not found"})
		}
		if student.AdvisorID == nil || *student.AdvisorID != lecturerID {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "not your advisee"})
		}
	}

	if err := s.Repo.UpdateReferenceStatus(refID, "verified", &userID); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed update status"})
	}

	_ = s.Repo.AddHistory(refID, "submitted", "verified", userID, "")

	return c.JSON(fiber.Map{
		"message": "achievement verified",
		"status":  "verified",
	})
}

// ------------------------- REJECT ----------------------------
func (s *AchievementService) RejectAchievement(c *fiber.Ctx) error {
	refID := c.Params("id")
	if refID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "missing id"})
	}

	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	if role != "lecturer" && role != "admin" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "only lecturer or admin can reject"})
	}

	var body struct {
		Note string `json:"note"`
	}
	if err := c.BodyParser(&body); err != nil || body.Note == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "rejection note required"})
	}

	ref, err := s.Repo.GetReferenceByID(refID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "reference not found"})
	}

	if ref.Status != "submitted" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "only submitted achievements can be rejected"})
	}

	student, err := s.StudentRepo.GetByID(ref.StudentID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "student not found"})
	}

	if role == "lecturer" {
		lecturerID, err := s.StudentRepo.GetLecturerIDByUserID(userID)
		if err != nil {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "lecturer profile not found"})
		}
		if student.AdvisorID == nil || *student.AdvisorID != lecturerID {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "not your advisee"})
		}
	}

	if err := s.Repo.RejectReference(refID, userID, body.Note); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed reject achievement"})
	}

	_ = s.Repo.AddHistory(refID, "submitted", "rejected", userID, body.Note)

	_ = s.Repo.CreateNotification(
		ref.StudentID,
		"Prestasi Ditolak",
		"Prestasi Anda ditolak dengan catatan: "+body.Note,
	)

	return c.JSON(fiber.Map{
		"message": "achievement rejected",
		"status":  "rejected",
	})
}

// ------------------------- HISTORY ----------------------------
func (s *AchievementService) GetAchievementHistory(c *fiber.Ctx) error {
	refID := c.Params("id")
	if refID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "missing id"})
	}

	s.Repo.EnsureDBs()

	var rows []struct {
		ID               string    `db:"id" json:"id"`
		AchievementRefID string    `db:"achievement_ref_id" json:"achievement_ref_id"`
		OldStatus        string    `db:"old_status" json:"old_status"`
		NewStatus        string    `db:"new_status" json:"new_status"`
		ChangedBy        string    `db:"changed_by" json:"changed_by"`
		Note             string    `db:"note" json:"note"`
		ChangedAt        time.Time `db:"changed_at" json:"changed_at"`
	}

	q := `SELECT id, achievement_ref_id, old_status, new_status, changed_by, note, changed_at
		  FROM achievement_history WHERE achievement_ref_id=$1 ORDER BY changed_at ASC`

	if err := s.Repo.Psql.Select(&rows, q, refID); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":  "failed retrieve history",
			"detail": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"history": rows})
}

// ------------------------- ATTACHMENT ----------------------------
func (s *AchievementService) UploadAttachment(c *fiber.Ctx) error {
	refID := c.Params("id")
	if refID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "missing id",
		})
	}

	// user dari JWT
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing user in token",
		})
	}
	userIDStr := userID.(string)

	// ===============================
	// Ambil file dari form-data
	// ===============================
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "file is required",
		})
	}

	// pastikan folder uploads ada
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed create upload directory",
		})
	}

	// ===============================
	// Generate nama file aman & unik
	// ===============================
	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf(
		"%s_%s%s",
		refID,
		uuid.New().String(),
		ext,
	)

	filePath := filepath.Join(uploadDir, fileName)

	// ===============================
	// Simpan file ke folder uploads/
	// ===============================
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed save file",
		})
	}

	// ===============================
	// Tentukan file type
	// ===============================
	fileType := c.FormValue("file_type")
	if fileType == "" {
		fileType = strings.TrimPrefix(ext, ".")
	}

	// ===============================
	// Simpan ke database
	// ===============================
	s.Repo.EnsureDBs()

	if err := s.Repo.AddAttachment(refID, filePath, fileType, userIDStr); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":  "failed save attachment",
			"detail": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":  "attachment uploaded",
		"file_url": filePath,
	})
}


// ------------------------- DELETE ----------------------------
func (s *AchievementService) DeleteAchievement(c *fiber.Ctx) error {
	refID := c.Params("id")
	if refID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "missing id"})
	}

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "missing user in token"})
	}
	userIDStr := userID.(string)

	s.Repo.EnsureDBs()

	student, err := s.StudentRepo.FindByUserID(userIDStr)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "student profile not found"})
	}

	ref, err := s.Repo.GetReferenceByID(refID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "achievement not found"})
	}

	if ref.StudentID != student.ID {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "not allowed"})
	}

	if ref.Status != "draft" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "only draft achievements can be deleted"})
	}

	if err := s.Repo.SoftDeleteMongo(ref.MongoAchievementID); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":  "failed to soft delete mongo",
			"detail": err.Error(),
		})
	}

	if err := s.Repo.SoftDeleteReference(refID); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":  "failed to update reference",
			"detail": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "achievement deleted successfully",
		"id":      refID,
	})
}

// ------------------------- ADVISEE ----------------------------
func (s *AchievementService) GetAdviseeAchievements(c *fiber.Ctx) error {
	lecturerID := c.Params("id")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	studentIDs, err := s.StudentRepo.GetStudentIDsByAdvisor(lecturerID)
	if err != nil || len(studentIDs) == 0 {
		return c.JSON(fiber.Map{
			"data":  []any{},
			"page":  page,
			"limit": limit,
		})
	}

	refs, err := s.Repo.GetReferencesByStudentIDs(studentIDs, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	results := []any{}
	for _, ref := range refs {
		detail, err := s.Repo.GetAchievementMongoDetail(ref.MongoAchievementID)
		if err != nil {
			fmt.Println("Mongo miss:", ref.MongoAchievementID)
			continue
		}

		results = append(results, fiber.Map{
			"reference": ref,
			"detail":    detail,
		})
	}

	return c.JSON(fiber.Map{
		"data":  results,
		"page":  page,
		"limit": limit,
	})
}

// GET /api/v1/achievements (ADMIN)
func (s *AchievementService) GetAll(c *fiber.Ctx) error {
	data, err := s.Repo.GetAllReferencesWithDetail()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed fetch achievements",
		})
	}
	return c.JSON(fiber.Map{
		"data": data,
	})
}

// PUT /api/v1/achievements/:id (STUDENT, DRAFT ONLY)
func (s *AchievementService) UpdateAchievement(c *fiber.Ctx) error {
	refID := c.Params("id")

	userID := c.Locals("user_id").(string)

	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid body",
		})
	}

	err := s.Repo.UpdateDraftAchievement(refID, userID, body)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "achievement updated",
	})
}
