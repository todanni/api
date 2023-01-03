package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/todanni/api/models"
)

type TaskRepository interface {
	CreateTask(task models.Task) (models.Task, error)
	GetTaskByID(taskID string) (models.Task, error)
	UpdateTask(task models.Task) (models.Task, error)
	UpdateTaskDone(taskID string, done bool) (models.Task, error)
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

func (r *taskRepo) UpdateTask(task models.Task) (models.Task, error) {
	result := r.db.Model(&task).Clauses(clause.Returning{}).Updates(map[string]interface{}{"title": task.Title, "description": task.Description, "done": task.Done, "assigned_to": task.AssignedTo, "deadline": task.Deadline})
	return task, result.Error
}

func (r *taskRepo) UpdateTaskDone(taskID string, done bool) (models.Task, error) {
	var task models.Task
	result := r.db.Model(&task).Where("id = ?", taskID).Clauses(clause.Returning{}).Updates(map[string]interface{}{"done": done})
	return task, result.Error
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepo{
		db: db,
	}
}
