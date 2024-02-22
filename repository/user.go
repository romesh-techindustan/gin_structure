package repository

import (
	"zucora/backend/models"

	"gorm.io/gorm"
)

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
    var user models.User
    err := r.db.First(&user,  "email = ?", id).Error
    return &user, err
}
