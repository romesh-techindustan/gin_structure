package controllers

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)
var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

func GenerateAccessToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 1).Unix(), // Token expiration time
	})

	return token.SignedString(jwtSecret)
	
}

func GenerateRefreshToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // Refresh token expiration time (7 days)
	})

	return token.SignedString(jwtSecret)
}