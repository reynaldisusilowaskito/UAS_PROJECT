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
		RoleID   string `json:"roleId"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	user := model.User{
		ID:           uuid.New().String(),
		Username:     body.Username,
		Email:        body.Email,
		PasswordHash: string(hash),
		FullName:     body.FullName,
		RoleID:       body.RoleID,
		IsActive:     true,
	}

	if err := s.UserRepo.Create(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(user)
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
