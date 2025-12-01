package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func Migrate(db *sqlx.DB) {

	fmt.Println("Running migration...")

	// ============================================
	// GENERATE HASH PASSWORD DEFAULT: "password123"
	// ============================================
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	password := string(hashedPassword)

	// ============================================
	// INSERT ROLES
	// ============================================
	db.Exec(`
		INSERT INTO roles (id, name, description)
		VALUES 
			(gen_random_uuid(), 'admin', 'Administrator'),
			(gen_random_uuid(), 'student', 'Mahasiswa'),
			(gen_random_uuid(), 'lecturer', 'Dosen Pembimbing')
		ON CONFLICT (name) DO NOTHING
	`)

	// Ambil role ID
	var adminRole, studentRole, lecturerRole string
	db.Get(&adminRole, "SELECT id FROM roles WHERE name='admin'")
	db.Get(&studentRole, "SELECT id FROM roles WHERE name='student'")
	db.Get(&lecturerRole, "SELECT id FROM roles WHERE name='lecturer'")

	// ============================================
	// INSERT PERMISSIONS
	// ============================================
	permissions := []struct {
		Resource string
		Action   string
	}{
		// Student
		{"achievements", "create"},
		{"achievements", "update"},
		{"achievements", "delete"},
		{"achievements", "submit"},
		{"students", "read-self"},

		// Lecturer
		{"achievements", "verify"},
		{"achievements", "reject"},
		{"students", "advisees"},

		// Admin
		{"achievements", "read-all"},
		{"achievements", "stats"},
		{"users", "create"},
		{"users", "read"},
		{"users", "update"},
		{"users", "delete"},
		{"users", "update-role"},
		{"students", "assign-advisor"},
	}

	permissionIDs := make(map[string]string)

	for _, p := range permissions {
		var permissionID string

		// Pakai _ agar tidak error
		_ = db.QueryRow(`
			INSERT INTO permissions (id, resource, action)
			VALUES (gen_random_uuid(), $1, $2)
			ON CONFLICT (resource, action) DO NOTHING
			RETURNING id
		`, p.Resource, p.Action).Scan(&permissionID)

		// Jika RETURNING id kosong, ambil id yang sudah ada
		if permissionID == "" {
			db.Get(&permissionID,
				"SELECT id FROM permissions WHERE resource=$1 AND action=$2",
				p.Resource, p.Action)
		}

		permissionIDs[p.Resource+":"+p.Action] = permissionID
	}

	// ============================================
	// INSERT ROLE_PERMISSIONS
	// ============================================
	assign := func(roleID string, key string) {
		permID := permissionIDs[key]
		_, err := db.Exec(`
			INSERT INTO role_permissions (id, role_id, permission_id)
			VALUES (gen_random_uuid(), $1, $2)
			ON CONFLICT DO NOTHING
		`, roleID, permID)

		if err != nil {
			log.Println("Error assign permission:", err)
		}
	}

	// Student
	for _, key := range []string{
		"achievements:create",
		"achievements:update",
		"achievements:delete",
		"achievements:submit",
		"students:read-self",
	} {
		assign(studentRole, key)
	}

	// Lecturer
	for _, key := range []string{
		"achievements:verify",
		"achievements:reject",
		"students:advisees",
	} {
		assign(lecturerRole, key)
	}

	// Admin
	for _, key := range []string{
		"achievements:read-all",
		"achievements:stats",
		"users:create",
		"users:read",
		"users:update",
		"users:delete",
		"users:update-role",
		"students:assign-advisor",
	} {
		assign(adminRole, key)
	}

	// ============================================
	// INSERT STUDENTS
	// ============================================
	students := []struct {
		Username string
		FullName string
		Email    string
		NIM      string
	}{
		{"student1", "Mahasiswa Satu", "student1@mail.com", "2023001"},
		{"student2", "Mahasiswa Dua", "student2@mail.com", "2023002"},
		{"student3", "Mahasiswa Tiga", "student3@mail.com", "2023003"},
		{"student4", "Mahasiswa Empat", "student4@mail.com", "2023004"},
		{"student5", "Mahasiswa Lima", "student5@mail.com", "2023005"},
		{"student6", "Mahasiswa Enam", "student6@mail.com", "2023006"},
	}

	for _, s := range students {

		var userID string

		err := db.QueryRow(`
			INSERT INTO users (id, username, email, password_hash, full_name, role_id)
			VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)
			RETURNING id
		`, s.Username, s.Email, password, s.FullName, studentRole).Scan(&userID)

		if err != nil {
			log.Println("Failed insert student:", err)
			continue
		}

		db.Exec(`
			INSERT INTO students (id, user_id, student_id, program_study, academic_year)
			VALUES (gen_random_uuid(), $1, $2, 'Informatika', '2023')
		`, userID, s.NIM)
	}

	// ============================================
	// INSERT LECTURERS
	// ============================================
	names := []string{"Tessa", "Arman", "Eto", "Alifian", "Indah", "Endah"}

	for _, name := range names {

		var userID string
		username := "lecturer_" + name
		email := username + "@mail.com"

		err := db.QueryRow(`
			INSERT INTO users (id, username, email, password_hash, full_name, role_id)
			VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)
			RETURNING id
		`, username, email, password, "Dosen "+name, lecturerRole).Scan(&userID)

		if err != nil {
			log.Println("Failed insert lecturer:", err)
			continue
		}

		db.Exec(`
			INSERT INTO lecturers (id, user_id, lecturer_id, department)
			VALUES (gen_random_uuid(), $1, $2, 'Teknik Informatika')
		`, userID, "D"+userID[:8])
	}

	// ============================================
	// INSERT ADMIN
	// ============================================
	db.Exec(`
		INSERT INTO users (id, username, email, password_hash, full_name, role_id)
		VALUES (gen_random_uuid(), 'admin', 'admin@mail.com', $1, 'Super Admin', $2)
		ON CONFLICT (username) DO NOTHING
	`, password, adminRole)

	fmt.Println("Migration completed!")
}
