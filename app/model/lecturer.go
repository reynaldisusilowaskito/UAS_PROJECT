package model

type Lecturer struct {
	ID       string `db:"id" json:"id"`
	UserID   string `db:"user_id" json:"user_id"`
	NIP      string `db:"nip" json:"nip"`
	FullName string `db:"full_name" json:"full_name"`
}
