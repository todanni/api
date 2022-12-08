package repository

import (
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

func (r *projectRepo) ListProjectsByUser(userID uint) ([]models.Project, error) {
	//TODO implement me
	panic("implement me")
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepo{
		db: db,
	}
}

func (r *projectRepo) CreateProject(project models.Project) (models.Project, error) {
	result := r.db.Create(&project)
	return project, result.Error
}
