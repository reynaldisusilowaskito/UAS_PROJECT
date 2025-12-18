package repository

import (
	"project_uas/app/model"
	"github.com/jmoiron/sqlx"
)

type LecturerRepo struct {
	DB *sqlx.DB
}

func NewLecturerRepo(db *sqlx.DB) *LecturerRepo {
	return &LecturerRepo{DB: db}
}

func (r *LecturerRepo) GetByUserID(userID string) (*model.Lecturer, error) {
	var lec model.Lecturer
	q := `
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
		WHERE user_id = $1
	`
	err := r.DB.Get(&lec, q, userID)
	return &lec, err
}


func (r *LecturerRepo) GetByID(id string) (*model.Lecturer, error) {
	lec := model.Lecturer{}
	q := `SELECT * FROM lecturers WHERE id=$1`
	err := r.DB.Get(&lec, q, id)
	return &lec, err
}

func (r *LecturerRepo) GetAll() ([]model.Lecturer, error) {
	var data []model.Lecturer
	q := `
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
	`
	err := r.DB.Select(&data, q)
	return data, err
}

