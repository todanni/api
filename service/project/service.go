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
	repo       repository.ProjectRepository
	middleware token.AuthMiddleware
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

func (s *projectService) ListProjectsHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)

	userID := accessToken.GetUserID()
	if userID == 0 {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}

	projects, err := s.repo.ListProjectsByUser(userID)
	if err != nil {
		http.Error(w, "couldn't retrieve projects", http.StatusInternalServerError)
		return
	}

	var response []ListProjectsResponse
	for _, project := range projects {
		response = append(response, ListProjectsResponse{
			ID:        project.ID,
			Name:      project.Name,
			Owner:     project.Owner,
			CreatedAt: project.CreatedAt,
			UpdatedAt: project.UpdatedAt,
			Members:   s.getMemberIDs(project.Members),
		})
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "couldn't marshall body", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(responseBody)
}

func (s *projectService) getMemberIDs(members []models.User) []uint {
	var ids []uint
	for _, member := range members {
		ids = append(ids, member.ID)
	}
	return ids
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

	if err = validation.ValidateStruct(&createRequest,
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

	response := CreateProjectResponse{
		ID:        project.ID,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
		Name:      project.Name,
		Owner:     project.Owner,
		Members:   s.getMemberIDs(project.Members),
	}

	responseBody, err := json.Marshal(response)
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

func (s *projectService) DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
