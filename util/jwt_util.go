package util

import (
	"fmt"
	"time"
	"github.com/dgrijalva/jwt-go"
)

const secretKey = "my_secret_key"

func GenerateToken(username string) (string, error) {
	// Creating a new token with the HS256 signing method and a MapClaims type
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,                           // Setting the username in the token claims
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Setting the expiration time of the token (24 hours)
	})

	// Signing the token with our secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("could not sign token: %v", err)
	}
	return tokenString, nil // Returning the signed token
}
