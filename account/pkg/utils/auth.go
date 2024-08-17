package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
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


// CreateJWTToken creates a JWT token with email, no_rekening, and expiration time in payload
func CreateJWTToken(jwtPayload dao.JWTField) (tokenString string, err error) {
	godotenv.Load("config.env")
	var JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
	// cast secret key to slice of bytes
	secretKey := []byte(JWT_SECRET_KEY)
	// Define the expiration time
	expirationTime := time.Now().Add(10 * time.Minute) // Example: 10 minute from now

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

	fmt.Println(JWT_SECRET_KEY,"============")

	return tokenString, nil
}

func ValidateJWTToken(tokenString string) (isValid bool, remark string, tokenData map[string]interface{}, err error) {
	godotenv.Load("config.env")
	var JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
	// cast secret key to slice of bytes
	secretKey := []byte(JWT_SECRET_KEY)
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
