package task

import "net/http"

const (
	APIPath = "/tasks"
)

func (s *taskService) routes() {
	r := s.router.PathPrefix(APIPath).Subrouter()

	r.HandleFunc("/", s.ListTasksHandler).Methods(http.MethodGet)
	r.HandleFunc("/", s.CreateTaskHandler).Methods(http.MethodPost)
	r.HandleFunc("/{id}", s.GetTaskHandler).Methods(http.MethodGet)
	r.HandleFunc("/{id}", s.UpdateTaskHandler).Methods(http.MethodPatch)
	r.HandleFunc("/{id}", s.DeleteTaskHandler).Methods(http.MethodDelete)
}
