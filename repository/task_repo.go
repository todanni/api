package repository

import (
	"gorm.io/gorm"

	"github.com/todanni/api/models"
)

type TaskRepository interface {
	GetTaskByID(taskID string) (models.Task, error)
	CreateTask(task models.Task) (models.Task, error)
	DeleteTask(taskID string) error
	ListTasksByUser(userID string) ([]models.Task, error)
	ListTasksByProject(projectID string) ([]models.Task, error)
}

type taskRepo struct {
	db *gorm.DB
}

func (r *taskRepo) GetTaskByID(taskID string) (models.Task, error) {
	var task models.Task
	result := r.db.First(&task, taskID)
	return task, result.Error
}

func (r *taskRepo) DeleteTask(taskID string) error {
	result := r.db.Delete(&models.Task{}, taskID)
	return result.Error
}

func (r *taskRepo) ListTasksByProject(projectID string) ([]models.Task, error) {
	var tasks []models.Task
	result := r.db.Where("project_id = ?", projectID).Find(&tasks)
	return tasks, result.Error
}

func (r *taskRepo) ListTasksByUser(userID string) ([]models.Task, error) {
	var tasks []models.Task
	result := r.db.Where("created_by = ?", userID).Or("assigned_to = ?", userID).Find(&tasks)
	return tasks, result.Error
}

func (r *taskRepo) CreateTask(task models.Task) (models.Task, error) {
	result := r.db.Create(&task)
	return task, result.Error
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepo{
		db: db,
	}
}
