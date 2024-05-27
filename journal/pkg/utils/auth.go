package utils

import (
	"golang.org/x/crypto/bcrypt"
)


func HashPassword(password string) (hashedPassword string, err error) {
	hashedPasswordByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	hashedPassword = string(hashedPasswordByte)
	return
}

func VerifyPassword(plainPassword, hashedPassword string)(err error) {
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return
}