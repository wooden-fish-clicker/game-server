package utils

import (
	"game-server/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Error hashing password:", err)
		return ""
	}

	return string(hashPassword)
}

func VerifyPassword(password string, hashPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err
}
