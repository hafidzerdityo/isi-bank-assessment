package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"hafidzresttemplate.com/dao"
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

// Hardcoded secret key
var secretKey = []byte("your_secret_key")

// CreateJWTToken creates a JWT token with email, no_rekening, and expiration time in payload
func CreateJWTToken(jwtPayload dao.JWTField) (tokenString string, err error) {
	// Define the expiration time
	expirationTime := time.Now().Add(10 * time.Minute) // Example: 1 day from now

	// Create the JWT claims, which includes the email, no_rekening, and exp
	claims := jwt.MapClaims{
		"no_rekening": jwtPayload.NoRekening,
		"exp":        expirationTime.Unix(),
	}

	// Create the JWT token with the claims and sign it with the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWTToken(tokenString string) (isValid bool, remark string, tokenData map[string]interface{}, err error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, jwt.ErrSignatureInvalid
		}
		return secretKey, nil
	})
	if err != nil {
		return false, fmt.Sprintf("failed to parse token, reason: %v", err.Error()), nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return false, "Your Token is Invalid", nil, jwt.ErrSignatureInvalid
	}

	// Token is valid and not expired
	return true, "Your Token is Valid", claims, nil
}
