package project

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/todanni/api/models"
	"github.com/todanni/api/repository"
	"github.com/todanni/api/token"
)

type ProjectsService interface {
	CreateProjectHandler(w http.ResponseWriter, r *http.Request)
	GetProjectHandler(w http.ResponseWriter, r *http.Request)
	UpdateProjectHandler(w http.ResponseWriter, r *http.Request)
	ListProjectsHandler(w http.ResponseWriter, r *http.Request)
	DeleteProjectHandler(w http.ResponseWriter, r *http.Request)
}

type projectService struct {
	router     *mux.Router
	middleware token.AuthMiddleware
	repo       repository.ProjectRepository
}

func (s *projectService) CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)

	userID := accessToken.GetUserID()
	if userID == 0 {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}

	var createRequest CreateProjectRequest
	err := json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = validation.ValidateStruct(createRequest,
		validation.Field(&createRequest.Name, validation.Required),
	); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	project, err := s.repo.CreateProject(models.Project{
		Name:  createRequest.Name,
		Owner: userID,
		Members: []models.User{
			{
				Model: gorm.Model{
					ID: userID,
				},
			},
		},
	})
	if err != nil {
		http.Error(w, "couldn't create project", http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(project)
	if err != nil {
		http.Error(w, "couldn't marshall body", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(responseBody)

}

func (s *projectService) GetProjectHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s *projectService) UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s *projectService) ListProjectsHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s *projectService) DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func NewProjectService(router *mux.Router, mw token.AuthMiddleware, repo repository.ProjectRepository) ProjectsService {
	service := &projectService{
		router:     router,
		repo:       repo,
		middleware: mw,
	}
	service.routes()
	return service
}
