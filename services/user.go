package services

import (
	"zucora/backend/models"
	"zucora/backend/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
    userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
    return &UserService{userRepository: userRepository}
}

func (s *UserService) Register(user *models.User) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password= string(hashedPassword)
    return s.userRepository.Create(user)
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
    return s.userRepository.GetByID(id)
}