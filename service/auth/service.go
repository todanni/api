package auth

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/goombaio/namegenerator"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	scopes "google.golang.org/api/oauth2/v2"
	"gorm.io/gorm"

	"github.com/todanni/api/config"
	"github.com/todanni/api/models"
	"github.com/todanni/api/repository"
	"github.com/todanni/api/token"
)

type AuthService interface {
	CallbackHandler(w http.ResponseWriter, r *http.Request)
	GetUserHandler(w http.ResponseWriter, r *http.Request)
}

type authService struct {
	router        *mux.Router
	userRepo      repository.UserRepository
	dashboardRepo repository.DashboardRepository
	projectRepo   repository.ProjectRepository
	middleware    token.AuthMiddleware
	config        config.Config
	oauthConfig   *oauth2.Config
}

func NewAuthService(
	router *mux.Router,
	cfg config.Config,
	userRepo repository.UserRepository,
	dashboardRepo repository.DashboardRepository,
	projectRepo repository.ProjectRepository,
	mw token.AuthMiddleware,
) AuthService {
	server := &authService{
		config:        cfg,
		router:        router,
		userRepo:      userRepo,
		dashboardRepo: dashboardRepo,
		projectRepo:   projectRepo,
		middleware:    mw,
	}
	server.routes()
	server.createOAuthConfig()
	return server
}

func (s *authService) createOAuthConfig() {
	decodedCredentials, err := b64.StdEncoding.DecodeString(s.config.GoogleCredentials)

	oauthConfig, err := google.ConfigFromJSON(decodedCredentials, scopes.OpenIDScope, scopes.UserinfoEmailScope, scopes.UserinfoProfileScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to oauthConfig: %v", err)
	}
	s.oauthConfig = oauthConfig
}

func (s *authService) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		log.Error(err)
		http.Error(w, "couldn't retrieve user", http.StatusInternalServerError)
		return
	}

	response := GetUserResponse{
		ID:          user.ID,
		DisplayName: user.DisplayName,
		ProfilePic:  user.ProfilePic,
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "couldn't marshall body", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(responseBody)
}

func (s *authService) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Received callback request")
	ctx := context.Background()

	code := r.URL.Query().Get("code")
	log.Info(s.oauthConfig)
	log.Info(s.oauthConfig.RedirectURL)

	tok, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		log.Errorf("Couldn't exchange keys: %v", err)
		http.Error(w, "couldn't exchange keys for code", http.StatusInternalServerError)
		return
	}

	// use the token to request the details about the user
	userInfo, err := s.getUserInfo(tok.AccessToken)
	if err != nil {
		log.Errorf("couldn't get user info from google: %v", err)
		http.Error(w, "couldn't get user info", http.StatusInternalServerError)
	}

	// Check if user exists
	userRecord, err := s.userRepo.GetUserByEmail(userInfo.Email)
	switch err {
	case gorm.ErrRecordNotFound:
		userRecord, err = s.userRepo.CreateUser(s.generateNewUserRecord(userInfo.Email, userInfo.ProfilePic))
		if err != nil {
			log.Errorf("Couldn't create user: %v", err)
			http.Error(w, "couldn't create new user", http.StatusInternalServerError)
			return
		}
	case nil:
		break
	default:
		log.Errorf("Couldn't check if user exists: %v", err)
		http.Error(w, "some error with user", http.StatusInternalServerError)
	}

	dashboards := make([]models.Dashboard, 0)
	projects := make([]models.Project, 0)

	if userRecord.ID != "" {
		dashboards, err = s.dashboardRepo.ListDashboardsByUser(userRecord.ID)
		if err != nil {
			log.Error("couldn't look up user dashboards")
		}

		projects, err = s.projectRepo.ListProjectsByUser(userRecord.ID)
		if err != nil {
			log.Error("couldn't look up user dashboards")
		}
	}

	accessToken := token.NewAccessToken()
	accessToken.SetUserID(userRecord.ID)
	accessToken.SetProjectsPermissions(projects)
	accessToken.SetDashboardPermissions(dashboards)

	signedToken, err := accessToken.SignToken([]byte(s.config.SigningKey))
	if err != nil {
		log.Errorf("Couldn't sign access token: %v", err)
		http.Error(w, "couldn't create access token", http.StatusInternalServerError)
		return
	}

	// Set access and refresh keys cookies
	http.SetCookie(w, &http.Cookie{
		Name:     token.AccessTokenCookieName,
		Value:    string(signedToken),
		Path:     "/",
		HttpOnly: true,
		Domain:   s.config.Domain,
	})

	w.Header().Set("Content-Type", "application/json")
	http.Redirect(w, r, s.config.RedirectURL, http.StatusFound)
}

type GoogleUserInfo struct {
	Email      string `json:"email"`
	ProfilePic string `json:"picture"`
}

func (s *authService) getUserInfo(accessToken string) (*GoogleUserInfo, error) {
	request, _ := http.NewRequest(http.MethodGet, "https://openidconnect.googleapis.com/v1/userinfo", nil)
	request.Header.Add("Authorization", "Bearer "+accessToken)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	var userInfo GoogleUserInfo
	err = json.NewDecoder(response.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func (s *authService) generateNewUserRecord(email, pic string) models.User {
	// generate display name
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	displayName := nameGenerator.Generate()

	// generate ID
	id := uuid.New().String()
	id = id[:8]

	return models.User{
		ID:          id,
		DisplayName: displayName,
		Email:       email,
		ProfilePic:  pic,
	}
}
