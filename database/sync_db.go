package database

import (
	"zucora/backend/models"
)

func SyncDb() {
	Db.AutoMigrate(&models.User{})
}
