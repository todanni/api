package repository

import (
	"gorm.io/gorm"

	"github.com/todanni/api/models"
)

type TaskRepository interface {
	CreateTask(task models.Task) (models.Task, error)
}

type taskRepo struct {
	db *gorm.DB
}

func (r *taskRepo) CreateTask(task models.Task) (models.Task, error) {
	result := r.db.Create(task)
	return task, result.Error
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepo{
		db: db,
	}
}
