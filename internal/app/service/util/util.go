package util

import (
	"akcom/internal/app/constants"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func ValidatePassword(password string, hashedPassword string) bool {
	// Comparing the password with the hash
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func GenerateHash(password string) (string, error) {
	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(hash), nil
}

func Int(v int) *int { return &v }

func IsRoleContain(str []constants.ROLE, sub constants.ROLE) bool {
	for _, v := range str {
		if v == sub {
			return true
		}
	}
	return false
}
