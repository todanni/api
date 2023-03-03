package routine

import (
	"time"
)

type CreateRoutineRequest struct {
	Name string `json:"name"`
	Days string `json:"days"`
}

type CreateRoutineResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Days      string    `json:"days"`
}

type ListRoutinesResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Days      string    `json:"days"`
}

type UpdateRoutineRequest struct {
	Name string `json:"name"`
	Days string `json:"days"`
}
