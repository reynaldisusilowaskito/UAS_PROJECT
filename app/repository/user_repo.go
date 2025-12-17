package repository

import (
	"project_uas/app/model"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	DB *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{DB: db}
}

// =====================
// READ
// =====================
func (r *UserRepo) GetAll() ([]model.User, error) {
	var users []model.User
	err := r.DB.Select(&users, `
		SELECT id, username, email, full_name, role_id, is_active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`)
	return users, err
}

func (r *UserRepo) GetByID(id string) (*model.User, error) {
	var user model.User
	err := r.DB.Get(&user, `
		SELECT id, username, email, full_name, role_id, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id)
	return &user, err
}

// =====================
// CREATE
// =====================
func (r *UserRepo) Create(user *model.User) error {
	_, err := r.DB.Exec(`
		INSERT INTO users 
		(id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,true,NOW(),NOW())
	`,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.RoleID,
	)
	return err
}

// =====================
// UPDATE
// =====================
func (r *UserRepo) Update(id string, user *model.User) error {
	_, err := r.DB.Exec(`
		UPDATE users
		SET username=$1, email=$2, full_name=$3, updated_at=NOW()
		WHERE id=$4
	`,
		user.Username,
		user.Email,
		user.FullName,
		id,
	)
	return err
}

// =====================
// DELETE (SOFT DELETE)
// =====================
func (r *UserRepo) Delete(id string) error {
	_, err := r.DB.Exec(`
		UPDATE users
		SET is_active=false, updated_at=NOW()
		WHERE id=$1
	`, id)
	return err
}

// =====================
// UPDATE ROLE
// =====================
func (r *UserRepo) UpdateRole(id string, roleID string) error {
	_, err := r.DB.Exec(`
		UPDATE users
		SET role_id=$1, updated_at=NOW()
		WHERE id=$2
	`, roleID, id)
	return err
}
