package project

import "time"

type CreateProjectRequest struct {
	Name string `json:"name"`
}

type CreateProjectResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Owner     uint      `json:"owner"`
	Members   []uint    `json:"members"`
}

type ListProjectsResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Owner     uint      `json:"owner"`
	Members   []uint    `json:"members"`
}
