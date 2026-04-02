package helper

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashPassword), nil
}

func VerifyPassword(hashedPassword, password string) error {
	if hashedPassword == "" || password == "" {
		return errors.New("password cannot be empty")
	}
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
