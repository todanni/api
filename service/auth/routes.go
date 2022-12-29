package auth

import "net/http"

const (
	CallbackHandler = "/auth/callback"
	GetUserHandler  = "/user"
)

func (s *authService) routes() {
	s.router.HandleFunc(CallbackHandler, s.CallbackHandler)

	r := s.router.PathPrefix(GetUserHandler).Subrouter()
	r.Use(s.middleware.JwtMiddleware)
	r.HandleFunc("/{id}", s.GetUserHandler).Methods(http.MethodGet)
}
