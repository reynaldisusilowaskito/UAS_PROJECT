package repository

import (
	"context"
	"time"

	"project_uas/app/model"
	"project_uas/database"

	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementRepo struct {
	Psql  *sqlx.DB
	Mongo *mongo.Database
}

func NewAchievementRepo(psql *sqlx.DB, mongo *mongo.Database) *AchievementRepo {
	return &AchievementRepo{Psql: psql, Mongo: mongo}
}

/*
	Exported helper to ensure DB handles exist.
	We export it so services can call r.EnsureDBs().
*/
func (r *AchievementRepo) EnsureDBs() {
	if r.Psql == nil {
		r.Psql = database.PostgresDB
	}
	if r.Mongo == nil {
		r.Mongo = database.MongoDB
	}
}

/* --------------------- Mongo Methods --------------------- */

func (r *AchievementRepo) CreateAchievementMongo(data model.Achievement) (string, error) {
	r.EnsureDBs()

	collection := r.Mongo.Collection("achievements")
	now := time.Now()
	data.CreatedAt = now
	data.UpdatedAt = now

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return "", err
	}

	oid := res.InsertedID.(primitive.ObjectID).Hex()
	return oid, nil
}

func (r *AchievementRepo) GetAchievementMongo(hexID string) (*model.Achievement, error) {
	r.EnsureDBs()

	collection := r.Mongo.Collection("achievements")
	oid, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return nil, err
	}

	var a model.Achievement
	if err := collection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&a); err != nil {
		return nil, err
	}
	return &a, nil
}

// UpdateAchievementMongo sets fields using $set (update map should be plain fields)
func (r *AchievementRepo) UpdateAchievementMongo(hexID string, update bson.M) error {
	r.EnsureDBs()

	collection := r.Mongo.Collection("achievements")
	oid, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return err
	}
	update["updated_at"] = time.Now()

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": oid}, bson.M{"$set": update})
	return err
}

// Helper to push array element(s) into files array
func (r *AchievementRepo) PushFileToAchievement(hexID string, files []string) error {
	r.EnsureDBs()

	collection := r.Mongo.Collection("achievements")
	oid, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return err
	}
	push := bson.M{"$push": bson.M{"files": bson.M{"$each": files}, "updated_at": time.Now()}}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": oid}, push)
	return err
}

/* -------------------- Postgres Methods ------------------- */

func (r *AchievementRepo) CreateReference(ref model.AchievementReference) error {
	r.EnsureDBs()

	q := `
		INSERT INTO achievement_references
		(id, student_id, mongo_achievement_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.Psql.Exec(q,
		ref.ID,
		ref.StudentID,
		ref.MongoAchievementID,
		ref.Status,
		ref.CreatedAt,
		ref.UpdatedAt,
	)
	return err
}

func (r *AchievementRepo) GetReferenceByID(id string) (*model.AchievementReference, error) {
	r.EnsureDBs()

	var ref model.AchievementReference
	q := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE id = $1
	`
	if err := r.Psql.Get(&ref, q, id); err != nil {
		return nil, err
	}
	return &ref, nil
}

func (r *AchievementRepo) GetReferenceByMongoID(mongoHex string) (*model.AchievementReference, error) {
	r.EnsureDBs()

	var ref model.AchievementReference
	q := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE mongo_achievement_id = $1
	`
	if err := r.Psql.Get(&ref, q, mongoHex); err != nil {
		return nil, err
	}
	return &ref, nil
}

func (r *AchievementRepo) GetAllForStudent(studentID string) ([]model.AchievementReference, error) {
	r.EnsureDBs()

	var data []model.AchievementReference
	q := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE student_id = $1
		ORDER BY created_at DESC
	`
	err := r.Psql.Select(&data, q, studentID)
	return data, err
}

func (r *AchievementRepo) UpdateReferenceStatus(refID string, status string, by *string) error {
	r.EnsureDBs()

	now := time.Now()

	switch status {
	case "submitted":
		_, err := r.Psql.Exec(`
				UPDATE achievement_references SET status=$1, submitted_at=$2, updated_at=$3 WHERE id=$4
			`, status, now, now, refID)
		return err

	case "verified":
		if by == nil {
			return nil
		}
		_, err := r.Psql.Exec(`
				UPDATE achievement_references SET status=$1, verified_at=$2, verified_by=$3, updated_at=$4 WHERE id=$5
			`, status, now, *by, now, refID)
		return err

	case "rejected":
		if by == nil {
			return nil
		}
		_, err := r.Psql.Exec(`
				UPDATE achievement_references SET status=$1, verified_at=$2, verified_by=$3, updated_at=$4 WHERE id=$5
			`, status, now, *by, now, refID)
		return err

	default:
		_, err := r.Psql.Exec(`UPDATE achievement_references SET status=$1, updated_at=$2 WHERE id=$3`, status, now, refID)
		return err
	}
}

func (r *AchievementRepo) AddHistory(refID, oldStatus, newStatus, changedBy, note string) error {
	r.EnsureDBs()

	_, err := r.Psql.Exec(`
		INSERT INTO achievement_history (id, achievement_ref_id, old_status, new_status, changed_by, note, changed_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6)
	`, refID, oldStatus, newStatus, changedBy, note, time.Now())
	return err
}

func (r *AchievementRepo) AddAttachment(refID, fileURL, fileType, uploadedBy string) error {
	r.EnsureDBs()

	_, err := r.Psql.Exec(`
		INSERT INTO achievement_attachments (id, achievement_ref_id, file_url, file_type, uploaded_by, uploaded_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)
	`, refID, fileURL, fileType, uploadedBy, time.Now())
	return err
}
