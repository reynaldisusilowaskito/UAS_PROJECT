package model

import "time"

type Student struct {
    ID           string    `db:"id" json:"id"`
    UserID       string    `db:"user_id" json:"user_id"`
    StudentID    string    `db:"student_id" json:"student_id"`
    ProgramStudy string    `db:"program_study" json:"program_study"`
    AcademicYear string    `db:"academic_year" json:"academic_year"`
    AdvisorID    *string   `db:"advisor_id" json:"advisor_id"`
    CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
