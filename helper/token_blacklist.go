package helper

import (
	"project_uas/database"
	"project_uas/app/repository"
)

// IsTokenBlacklisted mengecek apakah token sudah direvoke
func IsTokenBlacklisted(token string) bool {
	repo := repository.NewAuthRepo(database.PostgresDB)

	isRevoked, err := repo.IsTokenRevoked(token)
	if err != nil {
		// kalau DB error → lebih aman BLOCK
		return true
	}

	return isRevoked
}
