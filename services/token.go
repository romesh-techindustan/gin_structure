package services

import (
	"fmt"
	"net/http"
	"time"

	"zucora/backend/database"
	"zucora/backend/models"

	"github.com/golang-jwt/jwt/v5"
)

type TokenParams struct {
	TokenExpiration        int32
	RefreshTokenExpiration int32
	JWTSecret              []byte
	UserID                 string
	AuthToken              string
}

type IToken interface {
	GenerateAccessToken() (string, error)
	GenerateRefreshToken() (string, error)
	VerifyAccessToken() (*models.User, int, error) // returns (user id, http code, error)
}

func GenerateAccessToken(TP TokenParams) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": TP.UserID,
		"exp":    time.Now().Add(time.Hour * time.Duration(TP.TokenExpiration)).Unix(),
	})
	return token.SignedString(TP.JWTSecret)

}

func GenerateRefreshToken(TP TokenParams) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * time.Duration(TP.RefreshTokenExpiration)).Unix(),
	})
	return token.SignedString(TP.JWTSecret)
}

func VerifyAccessToken(TP TokenParams) (*models.User, int, error) {
	var user models.User
	token, err := jwt.Parse(TP.AuthToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["sub"])
		}
		return TP.JWTSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, http.StatusUnauthorized, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return nil, http.StatusUnauthorized, fmt.Errorf("token expired")
		}
		userID := claims["userID"].(string)
		result := database.Db.Where(&models.User{ID: userID}).Find(&user)
		if result.RowsAffected == 0 {
			return nil, http.StatusUnauthorized, fmt.Errorf("user not found")
		}
	}
	return &user, http.StatusOK, nil

}
