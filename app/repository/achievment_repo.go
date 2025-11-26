package repository

import (
	"context"
	"time"

	"project_uas/app/model"

	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive" 
)

type AchievementRepo struct {
	Psql *sqlx.DB
	Mongo *mongo.Database
}

func NewAchievementRepo(psql *sqlx.DB, mongo *mongo.Database) *AchievementRepo {
	return &AchievementRepo{Psql: psql, Mongo: mongo}
}

// -------------------------
// MongoDB
// -------------------------
func (r *AchievementRepo) CreateAchievementMongo(data model.Achievement) (string, error) {
	collection := r.Mongo.Collection("achievements")

	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return "", err
	}

	// convert ObjectID to hex string
	oid := res.InsertedID.(primitive.ObjectID).Hex()
	return oid, nil
}

func (r *AchievementRepo) GetAchievementMongo(id string) (*model.Achievement, error) {
	collection := r.Mongo.Collection("achievements")

	var doc model.Achievement
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&doc)

	return &doc, err
}

// -------------------------
// PostgreSQL Reference
// -------------------------
func (r *AchievementRepo) CreateReference(ref model.AchievementReference) error {
	q := `
		INSERT INTO achievement_references 
		(id, student_id, mongo_achievement_id, status)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.Psql.Exec(q,
		ref.ID,
		ref.StudentID,
		ref.MongoAchievementID,
		ref.Status,
	)

	return err
}

func (r *AchievementRepo) GetAllForStudent(studentID string) ([]model.AchievementReference, error) {
	var data []model.AchievementReference
	q := `SELECT * FROM achievement_references WHERE student_id=$1`
	err := r.Psql.Select(&data, q, studentID)
	return data, err
}
