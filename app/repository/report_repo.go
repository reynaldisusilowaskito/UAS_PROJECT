package repository

import (
	"github.com/jmoiron/sqlx"
)

type ReportRepo struct {
	DB *sqlx.DB
}

func NewReportRepo(db *sqlx.DB) *ReportRepo {
	return &ReportRepo{DB: db}
}

func (r *ReportRepo) CountAchievementsByStatus() (map[string]int, error) {
	rows, err := r.DB.Queryx(`
		SELECT status, COUNT(*) 
		FROM achievement_references
		GROUP BY status
	`)
	if err != nil {
		return nil, err
	}

	result := map[string]int{}
	for rows.Next() {
		var status string
		var total int
		rows.Scan(&status, &total)
		result[status] = total
	}
	return result, nil
}

func (r *ReportRepo) GetStudentAchievementReport(studentID string) (map[string]int, error) {

	rows, err := r.DB.Queryx(`
		SELECT status, COUNT(*) 
		FROM achievement_references
		WHERE student_id = $1
		GROUP BY status
	`, studentID)

	if err != nil {
		return nil, err
	}

	result := map[string]int{}
	for rows.Next() {
		var status string
		var total int
		rows.Scan(&status, &total)
		result[status] = total
	}

	return result, nil
}
