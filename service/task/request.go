package task

import (
	"time"
)

type CreateTaskRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Done        bool      `json:"done"`
	Deadline    time.Time `json:"deadline"`
	Project     uint      `json:"project"`
	CreatedBy   uint      `json:"created_by"`
	AssignedTo  uint      `json:"assigned_to"`
}
