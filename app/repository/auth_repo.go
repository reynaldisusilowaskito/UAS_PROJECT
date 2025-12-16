package repository

import (
	"time"

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



func (r *AuthRepo) RevokeToken(token string, userID string, exp time.Time) error {
	_, err := r.DB.Exec(`
		INSERT INTO revoked_tokens (token, user_id, expired_at)
		VALUES ($1, $2, $3)
	`, token, userID, exp)
	return err
}

func (r *AuthRepo) IsTokenRevoked(token string) (bool, error) {
	var exists bool
	err := r.DB.Get(&exists, `
		SELECT EXISTS (
			SELECT 1 FROM revoked_tokens WHERE token = $1
		)
	`, token)
	return exists, err
}