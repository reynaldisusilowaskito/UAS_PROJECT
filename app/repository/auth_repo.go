package repository

import (
	"project_uas/app/model"

	"github.com/jmoiron/sqlx"
)

type AuthRepo struct {
	DB *sqlx.DB
}

func NewAuthRepo(db *sqlx.DB) *AuthRepo {
	return &AuthRepo{DB: db}
}

// Ambil user berdasarkan username
func (r *AuthRepo) FindByUsername(username string) (*model.User, error) {
	user := model.User{}
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	err := r.DB.Get(&user, query, username)
	return &user, err
}

// Ambil role name berdasarkan role_id (Admin, Mahasiswa, Dosen Wali)
func (r *AuthRepo) GetRoleNameByID(roleID string) (string, error) {
	var role string
	query := `SELECT name FROM roles WHERE id = $1`
	err := r.DB.Get(&role, query, roleID)
	return role, err
}

// Ambil permissions berdasarkan role
func (r *AuthRepo) GetPermissionsByRole(roleID string) ([]string, error) {
	var perms []string
	query := `
		SELECT CONCAT(resource, ':', action)
		FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = $1
	`
	err := r.DB.Select(&perms, query, roleID)
	return perms, err
}
