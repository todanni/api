package project

import (
	"encoding/json"
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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

	ListProjectMembers(w http.ResponseWriter, r *http.Request)
	AddProjectMember(w http.ResponseWriter, r *http.Request)
	RemoveProjectMember(w http.ResponseWriter, r *http.Request)
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
	// Get the project ID from the request
	params := mux.Vars(r)
	projectID := params["id"]

	projectIDStr, err := strconv.ParseUint(projectID, 10, 32)
	if err != nil {
		log.Error(err)
		http.Error(w, "invalid project ID", http.StatusBadRequest)
		return
	}

	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)
	if !accessToken.HasProjectPermission(uint(projectIDStr)) {
		http.Error(w, "you don't have access to this project", http.StatusForbidden)
		return
	}

	project, err := s.repo.GetProjectByID(projectID)
	if err != nil {
		log.Error(err)
		http.Error(w, "couldn't find project", http.StatusNotFound)
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

func (s *projectService) DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectID := params["id"]

	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)
	userID := accessToken.GetUserID()
	if userID == 0 {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}

	project, err := s.repo.GetProjectByID(projectID)
	if err != nil {
		log.Error(err)
		http.Error(w, "couldn't find project", http.StatusNotFound)
		return
	}

	if project.Owner != userID {
		http.Error(w, "only the project owner can delete a project", http.StatusForbidden)
		return
	}

	err = s.repo.DeleteProject(projectID)
	if err != nil {
		log.Error(err)
		http.Error(w, "couldn't delete project", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *projectService) ListProjectMembers(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s *projectService) AddProjectMember(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s *projectService) RemoveProjectMember(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s *projectService) UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
