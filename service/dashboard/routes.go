package dashboard

import "net/http"

const (
	APIPath = "/dashboards"
)

func (s *dashboardService) routes() {
	r := s.router.PathPrefix(APIPath).Subrouter()

	r.HandleFunc("/", s.ListDashboardsHandler).Methods(http.MethodGet)
	r.HandleFunc("/", s.CreateDashboardHandler).Methods(http.MethodPost)
	r.HandleFunc("/{id}", s.GetDashboardHandler).Methods(http.MethodGet)
	r.HandleFunc("/{id}", s.UpdateDashboardHandler).Methods(http.MethodPatch)
	r.HandleFunc("/{id}", s.DeleteDashboardHandler).Methods(http.MethodDelete)
}
