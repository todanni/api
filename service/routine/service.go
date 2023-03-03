package routine

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

type RoutinesService interface {
	CreateRoutineHandler(w http.ResponseWriter, r *http.Request)
	GetRoutineHandler(w http.ResponseWriter, r *http.Request)
	UpdateRoutineHandler(w http.ResponseWriter, r *http.Request)
	ListRoutinesHandler(w http.ResponseWriter, r *http.Request)
	DeleteRoutineHandler(w http.ResponseWriter, r *http.Request)

	SaveRoutineRecordHandler(w http.ResponseWriter, r *http.Request)
	DeleteRoutineRecordHandler(w http.ResponseWriter, r *http.Request)
}

type routineService struct {
	router     *mux.Router
	repo       repository.RoutineRepository
	middleware token.AuthMiddleware
}

func (s *routineService) SaveRoutineRecordHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)

	userID := accessToken.GetUserID()
	if userID == "" {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}
}

func (s *routineService) DeleteRoutineRecordHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)

	userID := accessToken.GetUserID()
	if userID == "" {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}
}

func NewRoutineService(router *mux.Router, mw token.AuthMiddleware, repo repository.RoutineRepository) RoutinesService {
	service := &routineService{
		router:     router,
		repo:       repo,
		middleware: mw,
	}
	service.routes()
	return service
}

func (s *routineService) ListRoutinesHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)

	userID := accessToken.GetUserID()
	if userID == "" {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}

	routines, err := s.repo.ListRoutinesByUser(userID)
	if err != nil {
		http.Error(w, "couldn't retrieve routines", http.StatusInternalServerError)
		return
	}

	var response []ListRoutinesResponse
	for _, routine := range routines {
		response = append(response, ListRoutinesResponse{
			ID:        routine.ID,
			Name:      routine.Name,
			Days:      routine.Days,
			CreatedAt: routine.CreatedAt,
			UpdatedAt: routine.UpdatedAt,
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

func (s *routineService) CreateRoutineHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)

	userID := accessToken.GetUserID()
	if userID == "" {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}

	var createRequest CreateRoutineRequest
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

	routine, err := s.repo.CreateRoutine(models.Routine{
		Name: createRequest.Name,
		Days: createRequest.Days,
	})
	if err != nil {
		http.Error(w, "couldn't create routine", http.StatusInternalServerError)
		return
	}

	response := CreateRoutineResponse{
		ID:        routine.ID,
		CreatedAt: routine.CreatedAt,
		UpdatedAt: routine.UpdatedAt,
		Name:      routine.Name,
		Days:      createRequest.Days,
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "couldn't marshall body", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(responseBody)
}

func (s *routineService) GetRoutineHandler(w http.ResponseWriter, r *http.Request) {
	// Get the routine ID from the request
	params := mux.Vars(r)
	routineID := params["id"]

	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)
	userID := accessToken.GetUserID()
	if userID == "" {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}

	routine, err := s.repo.GetRoutineByID(routineID)
	if err != nil {
		log.Error(err)
		http.Error(w, "couldn't find routine", http.StatusNotFound)
		return
	}

	// We don't want to reveal to users whether a routine exists
	// if they are not the user that created it
	if routine.UserID != userID {
		http.Error(w, "routine not found", http.StatusNotFound)
		return
	}

	responseBody, err := json.Marshal(routine)
	if err != nil {
		http.Error(w, "couldn't marshall body", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(responseBody)
}

func (s *routineService) DeleteRoutineHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	routineID := params["id"]

	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)
	userID := accessToken.GetUserID()
	if userID == "" {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}

	routine, err := s.repo.GetRoutineByID(routineID)
	if err != nil {
		log.Error(err)
		http.Error(w, "couldn't find routine", http.StatusNotFound)
		return
	}

	if routine.UserID != userID {
		http.Error(w, "only the routine owner can delete a routine", http.StatusForbidden)
		return
	}

	err = s.repo.DeleteRoutine(routineID)
	if err != nil {
		log.Error(err)
		http.Error(w, "couldn't delete routine", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *routineService) UpdateRoutineHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	routineIDStr := params["routine_id"]

	accessToken := r.Context().Value(token.AccessTokenContextKey).(*token.ToDanniToken)
	userID := accessToken.GetUserID()
	if userID == "" {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}

	routine, err := s.repo.GetRoutineByID(routineIDStr)
	if err != nil {
		log.Error(err)
		http.Error(w, "couldn't find routine", http.StatusNotFound)
		return
	}

	if routine.UserID != userID {
		http.Error(w, "only the routine owner can update a routine", http.StatusForbidden)
		return
	}

	var updateRequest UpdateRoutineRequest
	err = json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	routineID, err := strconv.ParseUint(routineIDStr, 10, 32)
	if err != nil {
		log.Error(err)
		http.Error(w, "invalid member ID", http.StatusBadRequest)
		return
	}

	updatedRoutine, err := s.repo.UpdateRoutine(models.Routine{
		Model: gorm.Model{
			ID: uint(routineID),
		},
		Name: updateRequest.Name,
		Days: updateRequest.Days,
	})

	if err != nil {
		log.Error(err)
		http.Error(w, "couldn't update routine", http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(updatedRoutine)
	if err != nil {
		http.Error(w, "couldn't marshall body", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(responseBody)
}
