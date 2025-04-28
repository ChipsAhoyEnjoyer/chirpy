package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (hashedPW string, err error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func CheckPasswordHash(hash, password string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
