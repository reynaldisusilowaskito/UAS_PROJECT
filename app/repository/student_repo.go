package repository

import (
	"project_uas/app/model"
	"github.com/jmoiron/sqlx"
)

type StudentRepo struct {
	DB *sqlx.DB
}

func NewStudentRepo(db *sqlx.DB) *StudentRepo {
	return &StudentRepo{DB: db}
}

func (r *StudentRepo) GetByUserID(userID string) (*model.Student, error) {
	student := model.Student{}
	q := `SELECT * FROM students WHERE user_id=$1`
	err := r.DB.Get(&student, q, userID)
	return &student, err
}

func (r *StudentRepo) GetByID(id string) (*model.Student, error) {
	student := model.Student{}
	q := `SELECT * FROM students WHERE id=$1`
	err := r.DB.Get(&student, q, id)
	return &student, err
}

func (r *StudentRepo) GetAll() ([]model.Student, error) {
	var data []model.Student
	q := `SELECT * FROM students`
	err := r.DB.Select(&data, q)
	return data, err
}
