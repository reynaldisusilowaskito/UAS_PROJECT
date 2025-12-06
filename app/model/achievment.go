package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Achievement struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Category    string             `bson:"category" json:"category"` // Kompetisi / Publikasi / Seminar
	Level       string             `bson:"level" json:"level"`       // lokal / nasional / internasional
	Organizer   string             `bson:"organizer" json:"organizer"`
	Location    string             `bson:"location,omitempty" json:"location,omitempty"`
	EventDate   *time.Time         `bson:"event_date,omitempty" json:"event_date,omitempty"`
	Score       int                `bson:"score,omitempty" json:"score,omitempty"`

	Files []string `bson:"files,omitempty" json:"files,omitempty"`

	CreatedBy string    `bson:"created_by" json:"created_by,omitempty"` // akan diisi server
	CreatedAt time.Time `bson:"created_at" json:"created_at,omitempty"` // akan diisi server
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at,omitempty"` // akan diisi server
}

type AchievementHistory struct {
	ID               string    `db:"id" json:"id"`
	AchievementRefID string    `db:"achievement_ref_id" json:"achievement_ref_id"`
	OldStatus        string    `db:"old_status" json:"old_status"`
	NewStatus        string    `db:"new_status" json:"new_status"`
	ChangedBy        string    `db:"changed_by" json:"changed_by"`
	Note             string    `db:"note" json:"note"`
	ChangedAt        time.Time `db:"changed_at" json:"changed_at"`
}
