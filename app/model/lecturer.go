package model

import "time"

type Lecturer struct {
	ID         string    `db:"id" json:"id"`
	UserID     string    `db:"user_id" json:"user_id"`
	LecturerID string    `db:"lecturer_id" json:"lecturer_id"`
	Department string    `db:"department" json:"department"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
