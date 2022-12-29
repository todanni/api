package repository

import (
	"gorm.io/gorm"

	"github.com/todanni/api/models"
)

type UserRepository interface {
	CreateUser(user models.User) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByID(id string) (models.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{
		db: db,
	}
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

func (r *userRepo) GetUserByID(id string) (models.User, error) {
	var user models.User
	result := r.db.Raw("SELECT * FROM users WHERE id = ?", id).Scan(&user)
	return user, result.Error
}

//var projects []models.Project
//result := r.db.Raw(
//"SELECT * FROM projects INNER JOIN user_projects up on projects.id = up.project_id WHERE user_id=?", userID).
//Scan(&projects)
//return projects, result.Error
