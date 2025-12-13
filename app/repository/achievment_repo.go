package repository

import (
	"context"
	"time"
	"fmt"

	"project_uas/app/model"
	"project_uas/database"

	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	 "github.com/lib/pq"
)

/* ============================================================
   REPOSITORY STRUCT
============================================================ */

type AchievementRepo struct {
	Psql  *sqlx.DB
	Mongo *mongo.Database
}

func NewAchievementRepo(psql *sqlx.DB, mongo *mongo.Database) *AchievementRepo {
	return &AchievementRepo{Psql: psql, Mongo: mongo}
}

// Ensure both DB connections exist
func (r *AchievementRepo) EnsureDBs() {
	if r.Psql == nil {
		r.Psql = database.PostgresDB
	}
	if r.Mongo == nil {
		r.Mongo = database.MongoDB
	}
}

/* ============================================================
   MONGO METHODS (achievement documents)
============================================================ */

// Create new achievement document in MongoDB
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

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Retrieve achievement in Mongo by Hex ID
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

// Update any field in achievement document
func (r *AchievementRepo) UpdateAchievementMongo(hexID string, update bson.M) error {
	r.EnsureDBs()

	collection := r.Mongo.Collection("achievements")

	oid, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return err
	}

	update["updated_at"] = time.Now()

	_, err = collection.UpdateOne(context.Background(),
		bson.M{"_id": oid},
		bson.M{"$set": update},
	)
	return err
}

// Append file URLs into the "files" array
func (r *AchievementRepo) PushFileToAchievement(hexID string, files []string) error {
	r.EnsureDBs()

	collection := r.Mongo.Collection("achievements")

	oid, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return err
	}

	push := bson.M{
		"$push": bson.M{
			"files": bson.M{
				"$each": files,
			},
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": oid}, push)
	return err
}

/* ============================================================
   POSTGRES METHODS (achievement references, history, notifications)
============================================================ */

// Create new reference entry linking Student ↔ Mongo Achievement
func (r *AchievementRepo) CreateReference(ref model.AchievementReference) error {
	r.EnsureDBs()

	query := `
		INSERT INTO achievement_references
		(id, student_id, mongo_achievement_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.Psql.Exec(query,
		ref.ID,
		ref.StudentID,
		ref.MongoAchievementID,
		ref.Status,
		ref.CreatedAt,
		ref.UpdatedAt,
	)
	return err
}

// Get single reference by UUID (primary key)
func (r *AchievementRepo) GetReferenceByID(id string) (*model.AchievementReference, error) {
    r.EnsureDBs()

    // DEBUG 1 — cek DB yang dipakai
    var dbName string
    if err := r.Psql.Get(&dbName, "SELECT current_database()"); err == nil {
        fmt.Println("DEBUG: Connected to database =", dbName)
    }

    // DEBUG 2 — cek apakah ID dikirim
    fmt.Println("DEBUG: GetReferenceByID called with ID =", id)

    // DEBUG 3 — cek referensi apakah ada
    var exists bool
    if err := r.Psql.Get(&exists, "SELECT EXISTS(SELECT 1 FROM achievement_references WHERE id = $1)", id); err == nil {
        fmt.Println("DEBUG: EXISTS =", exists)
    }

    // Query asli
    var ref model.AchievementReference
    query := `
        SELECT id, student_id, mongo_achievement_id, status,
               submitted_at, verified_at, verified_by, rejection_note,
               created_at, updated_at
        FROM achievement_references
        WHERE id = $1
    `

    err := r.Psql.Get(&ref, query, id)
    if err != nil {
        fmt.Println("DEBUG: ERROR at SELECT:", err)
        return nil, err
    }

    fmt.Println("DEBUG: REFERENCE FOUND =", ref.ID)

    return &ref, nil
}


// Get reference by Mongo Achievement ID
func (r *AchievementRepo) GetReferenceByMongoID(mongoHex string) (*model.AchievementReference, error) {
	r.EnsureDBs()

	var ref model.AchievementReference

	query := `
		SELECT id, student_id, mongo_achievement_id, status,
			   submitted_at, verified_at, verified_by, rejection_note,
			   created_at, updated_at
		FROM achievement_references
		WHERE mongo_achievement_id = $1
	`

	if err := r.Psql.Get(&ref, query, mongoHex); err != nil {
		return nil, err
	}

	return &ref, nil
}

// List all achievements belonging to one student
func (r *AchievementRepo) GetAllForStudent(studentID string) ([]model.AchievementReference, error) {
	r.EnsureDBs()

	var list []model.AchievementReference

	query := `
		SELECT id, student_id, mongo_achievement_id, status,
			   submitted_at, verified_at, verified_by, rejection_note,
			   created_at, updated_at
		FROM achievement_references
		WHERE student_id = $1
		ORDER BY created_at DESC
	`

	err := r.Psql.Select(&list, query, studentID)
	return list, err
}

// Generic status update (submitted / verified / rejected / custom)
func (r *AchievementRepo) UpdateReferenceStatus(refID string, status string, by *string) error {
	r.EnsureDBs()

	now := time.Now()

	switch status {
	case "submitted":
		_, err := r.Psql.Exec(`
			UPDATE achievement_references 
			SET status=$1, submitted_at=$2, updated_at=$3
			WHERE id=$4
		`, status, now, now, refID)
		return err

	case "verified", "rejected":
		if by == nil {
			return nil
		}
		_, err := r.Psql.Exec(`
			UPDATE achievement_references 
			SET status=$1, verified_at=$2, verified_by=$3, updated_at=$4
			WHERE id=$5
		`, status, now, *by, now, refID)
		return err

	default:
		_, err := r.Psql.Exec(`
			UPDATE achievement_references 
			SET status=$1, updated_at=$2 
			WHERE id=$3
		`, status, now, refID)
		return err
	}
}

// Add to achievement history table
func (r *AchievementRepo) AddHistory(refID, oldStatus, newStatus, changedBy, note string) error {
	r.EnsureDBs()

	_, err := r.Psql.Exec(`
		INSERT INTO achievement_history 
		(id, achievement_ref_id, old_status, new_status, changed_by, note, changed_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6)
	`, refID, oldStatus, newStatus, changedBy, note, time.Now())

	return err
}

// Save attachment record in Postgres
func (r *AchievementRepo) AddAttachment(refID, fileURL, fileType, uploadedBy string) error {
	r.EnsureDBs()

	_, err := r.Psql.Exec(`
		INSERT INTO achievement_attachments 
		(id, achievement_ref_id, file_url, file_type, uploaded_by, uploaded_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)
	`, refID, fileURL, fileType, uploadedBy, time.Now())

	return err
}

// Simple "submit" wrapper (rarely used — service layer uses UpdateReferenceStatus)
func (r *AchievementRepo) Submit(id string) error {
    r.EnsureDBs()

    res, err := r.Psql.Exec(`
        UPDATE achievement_references
        SET status = 'submitted', submitted_at = NOW(), updated_at = NOW()
        WHERE id = $1 AND status = 'draft'
    `, id)

    if err != nil {
        return err
    }

    rows, _ := res.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("no rows updated (invalid id or not draft)")
    }

    return nil
}


// Create notification record for advisor
func (r *AchievementRepo) CreateNotification(userID, title, message string) error {
	r.EnsureDBs()

	_, err := r.Psql.Exec(`
		INSERT INTO notifications (user_id, title, message, is_read, created_at)
		VALUES ($1, $2, $3, false, NOW())
	`, userID, title, message)

	return err
}


func (r *AchievementRepo) SoftDeleteMongo(hexID string) error {
    r.EnsureDBs()

    collection := r.Mongo.Collection("achievements")
    oid, err := primitive.ObjectIDFromHex(hexID)
    if err != nil {
        return err
    }

    _, err = collection.UpdateOne(
        context.Background(),
        bson.M{"_id": oid},
        bson.M{
            "$set": bson.M{
                "deleted_at": time.Now(),
                "updated_at": time.Now(),
            },
        },
    )
    return err
}

func (r *AchievementRepo) SoftDeleteReference(refID string) error {
    r.EnsureDBs()

    _, err := r.Psql.Exec(`
        UPDATE achievement_references
        SET status = 'deleted', updated_at = NOW()
        WHERE id = $1
    `, refID)

    return err
}

// Get references by student IDs (pagination)
func (r *AchievementRepo) GetReferencesByStudentIDs(
	studentIDs []string,
	limit int,
	offset int,
) ([]model.AchievementReference, error) {

	var refs []model.AchievementReference

	query := `
		SELECT id, student_id, mongo_achievement_id, status,
		       created_at, updated_at
		FROM achievement_references
		WHERE student_id = ANY($1)
		  AND status = 'submitted'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	err := r.Psql.Select(&refs, query, pq.Array(studentIDs), limit, offset)
	return refs, err
}

func (r *AchievementRepo) GetAchievementMongoDetail(id string) (bson.M, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var result bson.M
	err = r.Mongo.Collection("achievements").
		FindOne(context.Background(), bson.M{"_id": objID}).
		Decode(&result)

	return result, err
}
