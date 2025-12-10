package model

import "time"

// AchievementReference tracks the administrative metadata in Postgres
type AchievementReference struct {
    ID                 string     `db:"id" json:"id"`
    StudentID          string     `db:"student_id" json:"student_id"`
    MongoAchievementID string     `db:"mongo_achievement_id" json:"mongo_achievement_id"`
    Status             string     `db:"status" json:"status"`
    SubmittedAt        *time.Time `db:"submitted_at" json:"submitted_at"`
    VerifiedAt         *time.Time `db:"verified_at" json:"verified_at"`
    VerifiedBy         *string    `db:"verified_by" json:"verified_by"`   // ← INI YANG HARUS DIUBAH
    RejectionNote      *string    `db:"rejection_note" json:"rejection_note"`
    CreatedAt          time.Time  `db:"created_at" json:"created_at"`
    UpdatedAt          time.Time  `db:"updated_at" json:"updated_at"`
}

