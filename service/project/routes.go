package project

import "net/http"

const (
	APIPath = "/projects"
)

func (s *projectService) routes() {
	r := s.router.PathPrefix(APIPath).Subrouter()
	r.Use(s.middleware.JwtMiddleware)

	r.HandleFunc("/", s.ListProjectsHandler).Methods(http.MethodGet)
	r.HandleFunc("/", s.CreateProjectHandler).Methods(http.MethodPost)
	r.HandleFunc("/{id}", s.GetProjectHandler).Methods(http.MethodGet)
	r.HandleFunc("/{id}", s.UpdateProjectHandler).Methods(http.MethodPatch)
	r.HandleFunc("/{id}", s.DeleteProjectHandler).Methods(http.MethodDelete)
}
