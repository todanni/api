package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"github.com/todanni/api/config"
	"github.com/todanni/api/models"
	"github.com/todanni/api/repository"
	"github.com/todanni/api/token"
)

type AuthService interface {
	CallbackHandler(w http.ResponseWriter, r *http.Request)
	//UserInfoHandler(w http.ResponseWriter, r *http.Request)
}

type authService struct {
	router        *mux.Router
	userRepo      repository.UserRepository
	dashboardRepo repository.DashboardRepository
	projectRepo   repository.ProjectRepository
	config        config.Config
	oauthConfig   *oauth2.Config
}

func NewAuthService(
	router *mux.Router,
	conf config.Config,
	userRepo repository.UserRepository,
	dashboardRepo repository.DashboardRepository,
	projectRepo repository.ProjectRepository,
	oauthConfig *oauth2.Config,
) AuthService {
	server := &authService{
		oauthConfig:   oauthConfig,
		config:        conf,
		router:        router,
		userRepo:      userRepo,
		dashboardRepo: dashboardRepo,
		projectRepo:   projectRepo,
	}
	server.routes()
	return server
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

	//idToken := tok.Extra("id_token").(string)
	// use the token to request the details about the user
	userInfo, err := s.getUserInfo(tok.AccessToken)
	if err != nil {
		log.Errorf("couldn't get user info from google: %v", err)
		http.Error(w, "couldn't get user info", http.StatusInternalServerError)
	}

	// Check if user exists
	userRecord, err := s.userRepo.GetUserByEmail(userInfo.Email)
	if err != nil {
		log.Errorf("Couldn't check if user exists: %v", err)
		http.Error(w, "some error with user", http.StatusInternalServerError)
		return
	}

	// User doesn't exist, we have to create it
	if userRecord.ID == 0 {
		userRecord, err = s.userRepo.CreateUser(models.User{
			FirstName:  userInfo.FirstName,
			LastName:   userInfo.LastName,
			Email:      userInfo.Email,
			ProfilePic: userInfo.ProfilePic,
		})
		if err != nil {
			log.Errorf("Couldn't create user: %v", err)
			http.Error(w, "couldn't create new user", http.StatusInternalServerError)
			return
		}
	}

	dashboards := make([]models.Dashboard, 0)
	projects := make([]models.Project, 0)

	if userRecord.ID != 0 {
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

	signedToken, err := accessToken.SignToken(s.config.SigningKey)
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
	})

	w.Header().Set("Content-Type", "application/json")
	http.Redirect(w, r, "/tasks", http.StatusFound)
}

type GoogleUserInfo struct {
	Email      string `json:"email"`
	FirstName  string `json:"given_name"`
	LastName   string `json:"family_name"`
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
