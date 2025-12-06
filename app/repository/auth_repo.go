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

// =====================
//  FIND USER BY USERNAME
// =====================
func (r *AuthRepo) FindByUsername(username string) (*model.User, error) {
	var user model.User

	err := r.DB.Get(&user, `
		SELECT 
			id, 
			username, 
			email, 
			password_hash, 
			full_name, 
			role_id,
			is_active, 
			created_at, 
			updated_at
		FROM users
		WHERE username = $1
	`, username)

	return &user, err
}

// =====================
//   GET ROLE NAME
// =====================
func (r *AuthRepo) GetRoleNameByID(roleID string) (string, error) {
	var role string

	err := r.DB.Get(&role, `
		SELECT name FROM roles WHERE id = $1
	`, roleID)

	return role, err
}

// =====================
//   GET PERMISSIONS
// =====================
func (r *AuthRepo) GetPermissionsByRole(roleID string) ([]string, error) {
	var perms []string

	err := r.DB.Select(&perms, `
		SELECT CONCAT(resource, ':', action)
		FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = $1
	`, roleID)

	return perms, err
}
