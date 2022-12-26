package token

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/todanni/api/models"
)

type AuthMiddlewareTestSuite struct {
	suite.Suite
	token string
}

func dummyHandler(w http.ResponseWriter, r *http.Request) {}

func (s *AuthMiddlewareTestSuite) SetupSuite() {
	accessToken := NewAccessToken()
	id := uuid.New().String()
	id = id[:8]
	accessToken.SetUserID(id)

	// Set a token with no projects
	projects := make([]models.Project, 0)
	accessToken.SetProjectsPermissions(projects)

	// Set a token with no dashboards
	dashboards := make([]models.Dashboard, 0)
	accessToken.SetDashboardPermissions(dashboards)

	signedToken, err := accessToken.SignToken([]byte(signingKey))
	require.NoError(s.T(), err)
	require.NotNil(s.T(), signedToken)
	s.token = string(signedToken)
}

func (s *AuthMiddlewareTestSuite) Test_AccessToken_Good() {
	router := mux.NewRouter()

	mw := NewAuthMiddleware(signingKey)

	router.Use(mw.JwtMiddleware)
	router.HandleFunc("/", dummyHandler).Methods("GET")

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "Bearer "+s.token)

	router.ServeHTTP(rw, req)
	require.Equal(s.T(), 200, rw.Code)
}

func (s *AuthMiddlewareTestSuite) Test_NoAuthHeader_CookiePresent() {
	router := mux.NewRouter()

	mw := NewAuthMiddleware(signingKey)

	router.Use(mw.JwtMiddleware)
	router.HandleFunc("/", dummyHandler).Methods("GET")

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:     AccessTokenCookieName,
		Value:    s.token,
		Path:     "/",
		HttpOnly: true,
	})

	router.ServeHTTP(rw, req)
	require.Equal(s.T(), 200, rw.Code)
}

func TestAuthMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}
