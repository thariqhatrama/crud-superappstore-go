package main

import (
	"FinalTask/utils"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	hash, _ := utils.HashPassword("admin123")
	fmt.Println(hash)
}

// CheckPasswordHash membandingkan plaintext password dengan hash-nya
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
