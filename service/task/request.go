package task

import (
	"time"
)

type CreateTaskRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Done        bool      `json:"done"`
	Deadline    time.Time `json:"deadline"`
	ProjectID   uint      `json:"project_id"`
	CreatedBy   uint      `json:"created_by"`
	AssignedTo  uint      `json:"assigned_to"`
}

type UpdateTaskRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Done        bool      `json:"done"`
	AssignedTo  uint      `json:"assigned_to"`
	Deadline    time.Time `json:"deadline"`
}
