package project

type CreateProjectRequest struct {
	Name string `json:"name"`
}

type ListProjectsResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Owner   uint   `json:"owner"`
	Members []uint `json:"members"`
}
