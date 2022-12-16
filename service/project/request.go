package project

import (
	"time"
)

type CreateProjectRequest struct {
	Name string `json:"name"`
}

type CreateProjectResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Owner     string    `json:"owner"`
}

type ListProjectsResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Owner     string    `json:"owner"`
}

type UpdateProjectRequest struct {
	Name  string `json:"name"`
	Owner uint   `json:"owner"`
}

type ListProjectMembersResponse struct {
	ID          string `json:"id"`
	ProfilePic  string `json:"profile_pic"`
	DisplayName string `json:"display_name"`
}
