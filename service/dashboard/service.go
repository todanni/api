package dashboard

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/todanni/api/repository"
)

type DashboardsService interface {
	CreateDashboardHandler(w http.ResponseWriter, r *http.Request)
	GetDashboardHandler(w http.ResponseWriter, r *http.Request)
	UpdateDashboardHandler(w http.ResponseWriter, r *http.Request)
	ListDashboardsHandler(w http.ResponseWriter, r *http.Request)
	DeleteDashboardHandler(w http.ResponseWriter, r *http.Request)
}

type dashboardService struct {
	router *mux.Router
	repo   repository.DashboardRepository
}

func NewDashboardService(r *mux.Router, repo repository.DashboardRepository) DashboardsService {
	service := &dashboardService{
		router: r,
		repo:   repo,
	}
	service.routes()
	return service
}

func (s *dashboardService) CreateDashboardHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s *dashboardService) GetDashboardHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s *dashboardService) UpdateDashboardHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s *dashboardService) ListDashboardsHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s *dashboardService) DeleteDashboardHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
