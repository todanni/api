package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/todanni/api/models"
)

type ProjectRepository interface {
	CreateProject(project models.Project) (models.Project, error)
	UpdateProject(project models.Project) (models.Project, error)
	ListProjectsByUser(userID string) ([]models.Project, error)
	GetProjectByID(projectID string) (models.Project, error)
	DeleteProject(projectID string) error
	ListProjectMembers(projectID string) ([]models.User, error)
	AddProjectMember(userID string, prjID uint) error
	RemoveProjectMember(userID string, prjID uint) error
}

type projectRepo struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepo{
		db: db,
	}
}

func (r *projectRepo) ListProjectsByUser(userID string) ([]models.Project, error) {
	var projects []models.Project
	result := r.db.Raw(
		"SELECT * FROM projects INNER JOIN user_projects up on projects.id = up.project_id WHERE user_id=?", userID).
		Scan(&projects)
	return projects, result.Error
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

func (r *projectRepo) ListProjectMembers(projectID string) ([]models.User, error) {
	var projectMembers []models.User
	result := r.db.Raw("SELECT * FROM users INNER JOIN user_projects up on users.id = up.user_id WHERE project_id = ?", projectID).
		Scan(&projectMembers)
	return projectMembers, result.Error
}

func (r *projectRepo) AddProjectMember(userID string, projectID uint) error {
	return r.db.Model(&models.User{ID: userID}).
		Association("Projects").
		Append(&models.Project{
			Model: gorm.Model{ID: projectID}})
}

func (r *projectRepo) RemoveProjectMember(userID string, projectID uint) error {
	return r.db.Model(&models.User{ID: userID}).
		Association("Projects").
		Delete(&models.Project{
			Model: gorm.Model{ID: projectID}})
}

func (r *projectRepo) UpdateProject(project models.Project) (models.Project, error) {
	result := r.db.Model(&project).Clauses(clause.Returning{}).Updates(project)
	return project, result.Error
}
