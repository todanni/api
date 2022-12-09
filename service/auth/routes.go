package auth

const (
	CallbackHandler = "/auth/callback"
)

func (s *authService) routes() {
	s.router.HandleFunc(CallbackHandler, s.CallbackHandler)
}
