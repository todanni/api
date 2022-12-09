package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/todanni/api/models"
)

type ProjectRepository interface {
	CreateProject(project models.Project) (models.Project, error)
	ListProjectsByUser(userID uint) ([]models.Project, error)
}

type projectRepo struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepo{
		db: db,
	}
}

func (r *projectRepo) ListProjectsByUser(userID uint) ([]models.Project, error) {
	var Projects []models.Project
	var user models.User

	result := r.db.Model(&models.User{}).Preload("Projects.Members").First(&user, userID)
	if result.Error != nil {
		return Projects, errors.New("couldn't find Projects")
	}

	return user.Projects, nil
}

func (r *projectRepo) CreateProject(project models.Project) (models.Project, error) {
	result := r.db.Create(&project)
	return project, result.Error
}
