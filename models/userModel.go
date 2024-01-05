package models

import (
	"gorm.io/gorm"
)

// User model

type User struct {
	gorm.Model
	ID          string `gorm:"primaryKey;"`
	Name        string
	Email       string `gorm:"unique;"`
	Password    string
	OTP         string
	IsSuperuser bool `gorm:"default:false"`
}
