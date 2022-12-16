package project

import (
	"encoding/json"
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

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
	if userID == "" {
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

func (s *projectService) getMemberIDs(members []models.User) []string {
	var ids []string
	for _, member := range members {
		ids = append(ids, member.ID)
	}
	return ids
}

func (s *projectService) CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)

	userID := accessToken.GetUserID()
	if userID == "" {
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
				ID: userID,
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
	if userID == "" {
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

	projectMembers, err := s.repo.ListProjectMembers(projectID)
	if err != nil {
		log.Error(err)
		http.Error(w, "couldn't list project members", http.StatusInternalServerError)
		return
	}

	var response []ListProjectMembersResponse
	for _, member := range projectMembers {
		response = append(response, ListProjectMembersResponse{
			ID:          member.ID,
			ProfilePic:  member.ProfilePic,
			DisplayName: member.DisplayName,
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

func (s *projectService) AddProjectMember(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectIDStr := params["project_id"]
	memberID := params["member_id"]

	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)
	userID := accessToken.GetUserID()
	if userID == "" {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}

	project, err := s.repo.GetProjectByID(projectIDStr)
	if err != nil {
		log.Error(err)
		http.Error(w, "couldn't find project", http.StatusNotFound)
		return
	}

	if project.Owner != userID {
		http.Error(w, "only the project owner can add members to a project", http.StatusForbidden)
		return
	}

	if userID == memberID {
		http.Error(w, "you're already a part of this project", http.StatusBadRequest)
		return
	}

	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		log.Error(err)
		http.Error(w, "invalid member ID", http.StatusBadRequest)
		return
	}

	err = s.repo.AddProjectMember(memberID, uint(projectID))
	if err != nil {
		log.Error(err)
		http.Error(w, "couldn't add member to project", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *projectService) RemoveProjectMember(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectIDStr := params["project_id"]
	memberID := params["member_id"]

	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)
	userID := accessToken.GetUserID()
	if userID == "" {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}

	project, err := s.repo.GetProjectByID(projectIDStr)
	if err != nil {
		log.Error(err)
		http.Error(w, "couldn't find project", http.StatusNotFound)
		return
	}

	if project.Owner != userID {
		http.Error(w, "only the project owner can add members to a project", http.StatusForbidden)
		return
	}

	if userID == memberID {
		http.Error(w, "you're already a part of this project", http.StatusBadRequest)
		return
	}

	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		log.Error(err)
		http.Error(w, "invalid member ID", http.StatusBadRequest)
		return
	}

	err = s.repo.RemoveProjectMember(memberID, uint(projectID))
	if err != nil {
		log.Error(err)
		http.Error(w, "couldn't remove member from project", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *projectService) UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
