package repository

import (
	"project_uas/app/model"
	"github.com/jmoiron/sqlx"
)

type AdminRepo struct {
	DB *sqlx.DB
}

func NewAdminRepo(db *sqlx.DB) *AdminRepo {
	return &AdminRepo{DB: db}
}

func (r *AdminRepo) GetAllUsers() ([]model.User, error) {
	var users []model.User
	q := `SELECT * FROM users`
	err := r.DB.Select(&users, q)
	return users, err
}

func (r *AdminRepo) UpdateUserStatus(userID string, active bool) error {
	q := `UPDATE users SET is_active=$1 WHERE id=$2`
	_, err := r.DB.Exec(q, active, userID)
	return err
}
