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

func (r *StudentRepo) FindByUserID(userID string) (*model.Student, error) {
    var s model.Student

    err := r.DB.Get(&s, `
        SELECT
            id,
            user_id,
            student_id,
            program_study,
            academic_year,
            advisor_id,
            created_at
        FROM students
        WHERE user_id = $1
    `, userID)

    if err != nil {
        return nil, err
    }

    return &s, nil
}

func (r *StudentRepo) GetByID(id string) (*model.Student, error) {
    var s model.Student

    err := r.DB.Get(&s, `
        SELECT
            id,
            user_id,
            student_id,
            program_study,
            academic_year,
            advisor_id,
            created_at
        FROM students
        WHERE id = $1
    `, id)

    if err != nil {
        return nil, err
    }

    return &s, nil
}

func (r *StudentRepo) GetAll() ([]model.Student, error) {
    var data []model.Student

    err := r.DB.Select(&data, `
        SELECT
            id,
            user_id,
            student_id,
            program_study,
            academic_year,
            advisor_id,
            created_at
        FROM students
    `)

    return data, err
}
