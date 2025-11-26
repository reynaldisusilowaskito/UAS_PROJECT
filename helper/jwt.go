package helper

import (
	"time"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID string, username string, role string, permissions []string) (string, error) {

	claims := jwt.MapClaims{
		"user_id":     userID,
		"username":    username,
		"role":        role,
		"permissions": permissions,
		"exp":         time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(os.Getenv("JWT_SECRET"))

	return token.SignedString(secret)
}

func GenerateRefreshToken(userID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // refresh 7 hari
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(os.Getenv("JWT_SECRET"))

	return token.SignedString(secret)
}
