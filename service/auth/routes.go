package auth

const (
	CallbackHandler = "/auth/callback"
	//UserInfoHandler = "/auth/user-info"
)

func (s *authService) routes() {
	s.router.HandleFunc(CallbackHandler, s.CallbackHandler)

	// only UserInfoHandler requires auth
	//s.router.Handle(UserInfoHandler, token.NewAuthenticationMiddleware(s.UserInfoHandler)).Methods(http.MethodGet)
}
