package token

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AuthMiddlewareTestSuite struct {
	suite.Suite
	signingKey string
}

const (
	token = ""
)

func dummyHandler(w http.ResponseWriter, r *http.Request) {}

func (s *AuthMiddlewareTestSuite) SetupSuite() {
	s.signingKey = ""
}

func (s *AuthMiddlewareTestSuite) Test_AccessToken_Good() {
	router := mux.NewRouter()

	mw := NewAuthMiddleware(s.signingKey)

	router.Use(mw.JwtMiddleware)
	router.HandleFunc("/", dummyHandler).Methods("GET")

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "Bearer "+token)

	router.ServeHTTP(rw, req)
	require.Equal(s.T(), 200, rw.Code)

}

func TestAuthMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}
