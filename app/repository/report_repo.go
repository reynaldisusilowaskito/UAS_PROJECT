package repository

import (
	"project_uas/app/model"
	"github.com/jmoiron/sqlx"
)

type ReportRepo struct {
	DB *sqlx.DB
}

func NewReportRepo(db *sqlx.DB) *ReportRepo {
	return &ReportRepo{DB: db}
}

func (r *ReportRepo) CountAchievementsByStatus() ([]model.ReportItem, error) {
	var data []model.ReportItem
	q := `
		SELECT status AS name, COUNT(*) AS total
		FROM achievement_references
		GROUP BY status
	`
	err := r.DB.Select(&data, q)
	return data, err
}
