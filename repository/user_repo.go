package repository

import (
	"gorm.io/gorm"

	"github.com/todanni/api/models"
)

type UserRepository interface {
	CreateUser(user models.User) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{}
}

func (r *userRepo) CreateUser(user models.User) (models.User, error) {
	result := r.db.Create(&user)
	return user, result.Error
}

func (r *userRepo) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	return user, result.Error
}
