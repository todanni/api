package task

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"

	"github.com/todanni/api/models"
	"github.com/todanni/api/repository"
	"github.com/todanni/api/token"
)

type TasksService interface {
	CreateTaskHandler(w http.ResponseWriter, r *http.Request)
	GetTaskHandler(w http.ResponseWriter, r *http.Request)
	UpdateTaskHandler(w http.ResponseWriter, r *http.Request)
	ListTasksHandler(w http.ResponseWriter, r *http.Request)
	DeleteTaskHandler(w http.ResponseWriter, r *http.Request)
}

type taskService struct {
	router     *mux.Router
	middleware token.AuthMiddleware
	taskRepo   repository.TaskRepository
}

func NewTaskService(r *mux.Router, taskRepo repository.TaskRepository, mw token.AuthMiddleware) TasksService {
	service := &taskService{
		router:     r,
		taskRepo:   taskRepo,
		middleware: mw,
	}
	service.routes()
	return service
}

func (s *taskService) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Read the user's JWT and get the user ID from it
	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)

	userID := accessToken.GetUserID()
	if userID == "" {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}

	// Read body
	var createRequest CreateTaskRequest
	err := json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the user belongs to the specified in the request project
	if !accessToken.HasProjectPermission(createRequest.ProjectID) {
		http.Error(w, "user unauthorized for this project", http.StatusForbidden)
	}

	if err = validation.ValidateStruct(&createRequest,
		validation.Field(&createRequest.Title, validation.Required),
		validation.Field(&createRequest.ProjectID, validation.Required),
	); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call DB and persist task
	task, err := s.taskRepo.CreateTask(models.Task{
		Title:       createRequest.Title,
		Description: createRequest.Description,
		Done:        createRequest.Done,
		ProjectID:   createRequest.ProjectID,
		CreatedBy:   userID,
		AssignedTo:  createRequest.AssignedTo,
		Deadline:    createRequest.Deadline,
	})

	if err != nil {
		http.Error(w, "couldn't create task", http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "couldn't marshall body", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(responseBody)
}

func (s *taskService) ListTasksHandler(w http.ResponseWriter, r *http.Request) {
	// Read the user's JWT and get the user ID from it
	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)

	userID := accessToken.GetUserID()
	if userID == "" {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}

	tasks, err := s.taskRepo.ListTasksByUser(userID)
	if err != nil {
		http.Error(w, "couldn't look up tasks for user", http.StatusInternalServerError)
		return
	}
	responseBody, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, "couldn't marshall body", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(responseBody)
}

func (s *taskService) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s *taskService) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s *taskService) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
