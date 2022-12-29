package task

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
		log.Infof("user with ID %s doesn't have permissions for project %d",
			userID, createRequest.ProjectID)
		http.Error(w, "user unauthorized for this project", http.StatusForbidden)
		return
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

	projectID := r.URL.Query().Get("project_id")
	var tasks []models.Task
	var err error
	if projectID != "" {
		tasks, err = s.taskRepo.ListTasksByProject(projectID)
	} else {
		tasks, err = s.taskRepo.ListTasksByUser(userID)
	}

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
	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)

	userID := accessToken.GetUserID()
	if userID == "" {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}
	params := mux.Vars(r)
	taskID := params["id"]

	task, err := s.taskRepo.GetTaskByID(taskID)
	if err != nil {
		log.Error(err)
		http.Error(w, "couldn't look up task", http.StatusInternalServerError)
		return
	}

	if !accessToken.HasProjectPermission(task.ProjectID) {
		http.Error(w, "you don't have access", http.StatusForbidden)
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

func (s *taskService) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)

	userID := accessToken.GetUserID()
	if userID == "" {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}
	params := mux.Vars(r)
	taskID := params["id"]
	task, err := s.taskRepo.GetTaskByID(taskID)
	if err != nil {
		log.Error(err)
		http.Error(w, "couldn't look up task", http.StatusInternalServerError)
		return
	}

	if !accessToken.HasProjectPermission(task.ProjectID) {
		http.Error(w, "you don't have access", http.StatusForbidden)
		return
	}

	var updateRequest UpdateTaskRequest
	err = json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	taskIDUint, err := strconv.ParseUint(taskID, 10, 0)
	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	updatedTask, err := s.taskRepo.UpdateTask(models.Task{
		ID:          uint(taskIDUint),
		Title:       updateRequest.Title,
		Description: updateRequest.Description,
		Done:        updateRequest.Done,
		AssignedTo:  updateRequest.AssignedTo,
		Deadline:    updateRequest.Deadline,
	})

	responseBody, err := json.Marshal(updatedTask)
	if err != nil {
		http.Error(w, "couldn't marshall body", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(responseBody)
}

func (s *taskService) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)

	userID := accessToken.GetUserID()
	if userID == "" {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	taskID := params["id"]

	// Only the person who created the task can delete it
	task, err := s.taskRepo.GetTaskByID(taskID)
	if err != nil {
		log.Error(err)
		http.Error(w, "couldn't find task", http.StatusNotFound)
		return
	}

	if task.CreatedBy != userID {
		log.Error(err)
		http.Error(w, "only the person who created the task can delete it", http.StatusForbidden)
		return
	}

	err = s.taskRepo.DeleteTask(taskID)
	if err != nil {
		log.Error()
		http.Error(w, "couldn't delete task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
