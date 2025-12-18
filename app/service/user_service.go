package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"project_uas/app/model"
	"project_uas/app/repository"
)

type UserService struct {
	UserRepo *repository.UserRepo
}

func NewUserService(repo *repository.UserRepo) *UserService {
	return &UserService{UserRepo: repo}
}

// =====================
// GET ALL USERS
// =====================
func (s *UserService) GetAll(c *fiber.Ctx) error {
	users, err := s.UserRepo.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

// =====================
// GET USER BY ID
// =====================
func (s *UserService) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := s.UserRepo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}
	return c.JSON(user)
}

// =====================
// CREATE USER
// =====================
func (s *UserService) Create(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"fullName"`
		Role     string `json:"role"`

		// student
		StudentID    string `json:"student_id"`
		ProgramStudy string `json:"program_study"`
		AcademicYear string `json:"academic_year"`

		// lecturer
		LecturerID string `json:"lecturer_id"`
		Department string `json:"department"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	if body.Role != "student" && body.Role != "lecturer" && body.Role != "admin" {
		return c.Status(400).JSON(fiber.Map{"error": "invalid role"})
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	// ambil role_id dari table roles
	roleID, err := s.UserRepo.GetRoleIDByName(body.Role)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "role not found"})
	}

	userID := uuid.NewString()

	user := model.User{
		ID:           userID,
		Username:     body.Username,
		Email:        body.Email,
		PasswordHash: string(hash),
		FullName:     body.FullName,
		RoleID:       roleID,
		IsActive:     true,
	}

	tx, err := s.UserRepo.DB.Beginx()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "transaction failed"})
	}
	defer tx.Rollback()

	// =====================
	// INSERT USERS
	// =====================
	if err := s.UserRepo.CreateTx(tx, &user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// =====================
	// INSERT STUDENT
	// =====================
	if body.Role == "student" {
		if body.StudentID == "" {
			return c.Status(400).JSON(fiber.Map{"error": "student_id required"})
		}

		_, err := tx.Exec(`
			INSERT INTO students
			(id, user_id, student_id, program_study, academic_year)
			VALUES ($1,$2,$3,$4,$5)
		`,
			uuid.NewString(),
			userID,
			body.StudentID,
			body.ProgramStudy,
			body.AcademicYear,
		)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	}

	// =====================
	// INSERT LECTURER
	// =====================
	if body.Role == "lecturer" {
		if body.LecturerID == "" {
			return c.Status(400).JSON(fiber.Map{"error": "lecturer_id required"})
		}

		_, err := tx.Exec(`
			INSERT INTO lecturers
			(id, user_id, lecturer_id, department)
			VALUES ($1,$2,$3,$4)
		`,
			uuid.NewString(),
			userID,
			body.LecturerID,
			body.Department,
		)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	}

	if err := tx.Commit(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "commit failed"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "user created successfully",
		"user_id": userID,
		"role":    body.Role,
	})
}




// =====================
// UPDATE USER
// =====================
func (s *UserService) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		FullName string `json:"fullName"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	user := model.User{
		Username: body.Username,
		Email:    body.Email,
		FullName: body.FullName,
	}

	if err := s.UserRepo.Update(id, &user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "user updated"})
}

// =====================
// DELETE USER
// =====================
func (s *UserService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := s.UserRepo.Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "user deactivated"})
}

// =====================
// UPDATE ROLE
// =====================
func (s *UserService) UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")

	var body struct {
		RoleID string `json:"roleId"`
	}
	if err := c.BodyParser(&body); err != nil || body.RoleID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "roleId required"})
	}

	if err := s.UserRepo.UpdateRole(id, body.RoleID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "role updated"})
}
