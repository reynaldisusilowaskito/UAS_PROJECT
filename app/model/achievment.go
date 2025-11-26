package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Achievement struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID      string             `bson:"student_id" json:"student_id"`
	ReferenceID    string             `bson:"reference_id" json:"reference_id"`
	ProofFileURL   string             `bson:"proof_file_url" json:"proof_file_url"`
	Status         string             `bson:"status" json:"status"` // active | deleted
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}
