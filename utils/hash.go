// File: utils/hash.go
package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword menghasilkan hash bcrypt dari plaintext password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash membandingkan plaintext password dengan hash-nya
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
