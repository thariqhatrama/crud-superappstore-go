// File: utils/jwt.go
package utils

import (
	"time"

	"FinalTask/config"

	"github.com/golang-jwt/jwt/v4"
)

// GenerateJWT membuat token JWT dengan klaim user_id dan is_admin
func GenerateJWT(userID uint, isAdmin bool) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"is_admin": isAdmin,
		"exp":      time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JwtSecret)
}
