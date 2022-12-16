package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Done        bool           `json:"done"`
	Project     Project        `json:"-"`
	ProjectID   uint           `json:"project_id"`
	CreatedBy   string         `json:"created_by"`
	AssignedTo  uint           `json:"assigned_to"`
	Deadline    time.Time      `json:"deadline"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
