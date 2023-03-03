package routine

import "net/http"

const (
	APIPath = "/routines"
)

func (s *routineService) routes() {
	r := s.router.PathPrefix(APIPath).Subrouter()
	r.Use(s.middleware.JwtMiddleware)

	r.HandleFunc("/", s.ListRoutinesHandler).Methods(http.MethodGet)
	r.HandleFunc("/", s.CreateRoutineHandler).Methods(http.MethodPost)
	r.HandleFunc("/{id}", s.GetRoutineHandler).Methods(http.MethodGet)
	r.HandleFunc("/{id}", s.UpdateRoutineHandler).Methods(http.MethodPatch)
	r.HandleFunc("/{id}", s.DeleteRoutineHandler).Methods(http.MethodDelete)
}
