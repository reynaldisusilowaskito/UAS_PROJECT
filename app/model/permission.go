package model

type Permission struct {
	ID        string `db:"id" json:"id"`
	RoleID    string `db:"role_id" json:"role_id"`
	MenuName  string `db:"menu_name" json:"menu_name"`
	CanCreate bool   `db:"can_create" json:"can_create"`
	CanRead   bool   `db:"can_read" json:"can_read"`
	CanUpdate bool   `db:"can_update" json:"can_update"`
	CanDelete bool   `db:"can_delete" json:"can_delete"`
}
