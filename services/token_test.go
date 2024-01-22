package services

import (
	"testing"
)

func TestGenerateAccessToken(t *testing.T) {
	tp := &TokenParams{
		TokenExpiration: 1,
		UserID:          "testUserID",
	}
	_, err := GenerateAccessToken(*tp)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	tp := &TokenParams{
		RefreshTokenExpiration: 168,
		UserID:                 "testUserID",
	}
	_, err := GenerateRefreshToken(*tp)
	if err != nil {
		t.Errorf(err.Error())
	}
}
