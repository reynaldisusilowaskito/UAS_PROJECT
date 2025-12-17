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

// GetStudentIDsByAdvisor
func (r *StudentRepo) GetStudentIDsByAdvisor(advisorID string) ([]string, error) {

	var ids []string

	query := `
		SELECT id
		FROM students
		WHERE advisor_id = $1
	`

	err := r.DB.Select(&ids, query, advisorID)
	return ids, err
}

func (r *StudentRepo) GetLecturerIDByUserID(userID string) (string, error) {
	var lecturerID string
	err := r.DB.Get(&lecturerID, `
		SELECT id FROM lecturers WHERE user_id = $1
	`, userID)
	return lecturerID, err
}


// =====================
// GET STUDENT ACHIEVEMENTS
// =====================
func (r *StudentRepo) GetAchievementsByStudentID(studentID string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}

	query := `
		SELECT
			id,
			title,
			category,
			level,
			status,
			created_at
		FROM achievements
		WHERE student_id = $1
		ORDER BY created_at DESC
	`

	err := r.DB.Select(&data, query, studentID)
	return data, err
}

// =====================
// UPDATE ADVISOR
// =====================
func (r *StudentRepo) UpdateAdvisor(studentID string, advisorID string) error {
	_, err := r.DB.Exec(`
		UPDATE students
		SET advisor_id = $1
		WHERE id = $2
	`, advisorID, studentID)

	return err
}
