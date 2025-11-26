package model

type ReportItem struct {
	Name  string `db:"name" json:"name"`
	Total int    `db:"total" json:"total"`
}
