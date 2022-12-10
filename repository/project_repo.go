package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/todanni/api/models"
)

type ProjectRepository interface {
	CreateProject(project models.Project) (models.Project, error)
	UpdateProject(project models.Project) (models.Project, error)
	ListProjectsByUser(userID uint) ([]models.Project, error)
	GetProjectByID(projectID string) (models.Project, error)
	DeleteProject(projectID string) error
	GetProjectMembers(projectID string) error
	AddProjectMember(userID uint) error
	RemoveProjectMember(userID uint) error
}

type projectRepo struct {
	db *gorm.DB
}

func (r *projectRepo) UpdateProject(project models.Project) (models.Project, error) {
	// TODO: decide what we want to update this way
	return project, nil
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

func (r *projectRepo) DeleteProject(projectID string) error {
	result := r.db.Delete(&models.Project{}, projectID)
	return result.Error
}

func (r *projectRepo) GetProjectByID(projectID string) (models.Project, error) {
	var project models.Project
	result := r.db.First(&project, projectID)
	return project, result.Error
}

func (r *projectRepo) GetProjectMembers(projectID string) error {
	//TODO implement me
	panic("implement me")
}

func (r *projectRepo) AddProjectMember(userID uint) error {
	//TODO implement me
	panic("implement me")
}

func (r *projectRepo) RemoveProjectMember(userID uint) error {
	//TODO implement me
	panic("implement me")
}
