package service

import (
	"net/http"
	"log"
	"time"

	"github.com/gin-gonic/gin"
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
func (s *AchievementService) CreateAchievement(c *gin.Context) {
    var req model.Achievement

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }

    // ------------------------------
    // Ambil user dari JWT
    // ------------------------------
    userID := c.GetString("user_id")
    if userID == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user in token"})
        return
    }

    log.Println("DEBUG Token user_id =", userID)

    // ------------------------------
    // Ambil student profile berdasar user_id
    // ------------------------------
    student, err := s.StudentRepo.FindByUserID(userID)
    if err != nil {
        log.Println("DEBUG StudentRepo.FindByUserID error =", err)
        c.JSON(http.StatusNotFound, gin.H{
            "error":  "student profile not found",
            "detail": err.Error(),
        })
        return
    }

    log.Println("DEBUG FOUND STUDENT:", student.ID, student.UserID)

    // ------------------------------
    // Siapkan data untuk Mongo
    // ------------------------------
    req.CreatedBy = userID
    now := time.Now()

    req.CreatedAt = now
    req.UpdatedAt = now

    s.Repo.EnsureDBs()

    // ------------------------------
    // Insert MongoDB
    // ------------------------------
    mongoHex, err := s.Repo.CreateAchievementMongo(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":  "failed insert achievement into mongo",
            "detail": err.Error(),
        })
        return
    }

    // ------------------------------
    // Insert reference ke PostgreSQL
    // ------------------------------
    ref := model.AchievementReference{
        ID:                 uuid.New().String(),
        StudentID:          student.ID,
        MongoAchievementID: mongoHex,
        Status:             "draft",
        CreatedAt:          now,
        UpdatedAt:          now,
    }

    if err := s.Repo.CreateReference(ref); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":  "failed insert reference to postgres",
            "detail": err.Error(),
        })
        return
    }

    // Success
    c.JSON(http.StatusCreated, gin.H{
        "message": "achievement created",
        "data":    ref,
    })
}

// ------------------------- DETAIL ----------------------------
func (s *AchievementService) GetAchievementDetail(c *gin.Context) {
	refID := c.Param("id")
	if refID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}

	s.Repo.EnsureDBs()

	ref, err := s.Repo.GetReferenceByID(refID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reference not found"})
		return
	}

	ach, err := s.Repo.GetAchievementMongo(ref.MongoAchievementID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed retrieve achievement", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reference":   ref,
		"achievement": ach,
	})
}

// ----------------------- SUBMIT ----------------------------
func (s *AchievementService) SubmitAchievement(c *gin.Context) {
	refID := c.Param("id")
	if refID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user in token"})
		return
	}

	s.Repo.EnsureDBs()

	ref, err := s.Repo.GetReferenceByID(refID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reference not found"})
		return
	}

	// hanya owner (student) atau admin boleh submit
	role := c.GetString("role")
	if role != "admin" && ref.StudentID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed to submit this achievement"})
		return
	}

	// cek status saat ini
	if ref.Status == "submitted" || ref.Status == "verified" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "already submitted or processed"})
		return
	}

	if err := s.Repo.UpdateReferenceStatus(refID, "submitted", nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed update", "detail": err.Error()})
		return
	}

	_ = s.Repo.AddHistory(refID, ref.Status, "submitted", userID, "")

	c.JSON(http.StatusOK, gin.H{"message": "submitted"})
}

// ------------------------- VERIFY ----------------------------
func (s *AchievementService) VerifyAchievement(c *gin.Context) {
	refID := c.Param("id")
	if refID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user in token"})
		return
	}

	// only lecturer or admin can verify
	role := c.GetString("role")
	if role != "lecturer" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only lecturer or admin can verify"})
		return
	}

	s.Repo.EnsureDBs()

	ref, err := s.Repo.GetReferenceByID(refID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reference not found"})
		return
	}

	if ref.Status != "submitted" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only submitted achievements can be verified"})
		return
	}

	if err := s.Repo.UpdateReferenceStatus(refID, "verified", &userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed update status", "detail": err.Error()})
		return
	}

	_ = s.Repo.AddHistory(refID, "submitted", "verified", userID, "")

	c.JSON(http.StatusOK, gin.H{"message": "verified"})
}

// ------------------------- REJECT ----------------------------
func (s *AchievementService) RejectAchievement(c *gin.Context) {
	refID := c.Param("id")
	if refID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user in token"})
		return
	}

	// only lecturer or admin can reject
	role := c.GetString("role")
	if role != "lecturer" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only lecturer or admin can reject"})
		return
	}

	var body struct {
		Note string `json:"note"`
	}
	_ = c.ShouldBindJSON(&body)

	s.Repo.EnsureDBs()

	ref, err := s.Repo.GetReferenceByID(refID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reference not found"})
		return
	}

	if ref.Status != "submitted" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only submitted achievements can be rejected"})
		return
	}

	if err := s.Repo.UpdateReferenceStatus(refID, "rejected", &userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed update status", "detail": err.Error()})
		return
	}

	_ = s.Repo.AddHistory(refID, "submitted", "rejected", userID, body.Note)

	c.JSON(http.StatusOK, gin.H{"message": "rejected"})
}

// ------------------------- HISTORY ----------------------------
func (s *AchievementService) GetAchievementHistory(c *gin.Context) {
	refID := c.Param("id")
	if refID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed retrieve history", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": rows})
}

// ------------------------- ATTACHMENT ----------------------------
func (s *AchievementService) UploadAttachment(c *gin.Context) {
	refID := c.Param("id")
	if refID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user in token"})
		return
	}

	var body struct {
		FileURL  string `json:"file_url"`
		FileType string `json:"file_type"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.FileURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file_url required"})
		return
	}

	s.Repo.EnsureDBs()

	if err := s.Repo.AddAttachment(refID, body.FileURL, body.FileType, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed save attachment", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "attachment uploaded"})
}
