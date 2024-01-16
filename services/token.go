package services

import (
	"time"

	"zucora/backend/config"

	"github.com/golang-jwt/jwt/v5"
)

type TokenParams struct {
	Config     config.Config
	JWT_SECRET []byte
	USER_ID    string
}

type IToken interface {
	GenerateAccessToken() (string, error)
	GenerateRefreshToken() (string, error)
}

func GenerateAccessToken(TP TokenParams) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": TP.USER_ID,
		"exp":    time.Now().Add(time.Hour * time.Duration(TP.Config.TokenExpiration)).Unix(),
	})
	return token.SignedString(TP.JWT_SECRET)

}

func GenerateRefreshToken(TP TokenParams) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * time.Duration(TP.Config.RefreshTokenExpiration)).Unix(),
	})
	return token.SignedString(TP.JWT_SECRET)
}
