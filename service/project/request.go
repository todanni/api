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
	Owner     uint      `json:"owner"`
}

type ListProjectsResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Owner     uint      `json:"owner"`
}

type UpdateProjectRequest struct {
	Name  string `json:"name"`
	Owner uint   `json:"owner"`
}

type ListProjectMembersResponse struct {
	ID         uint   `json:"id"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ProfilePic string `json:"profile_pic"`
}
