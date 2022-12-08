package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/todanni/api/config"
	"github.com/todanni/api/database"
	"github.com/todanni/api/models"
	"github.com/todanni/api/repository"
	"github.com/todanni/api/service/project"
	"github.com/todanni/api/service/task"
	"github.com/todanni/api/token"
)

func main() {
	// Read config
	cfg, err := config.NewFromEnv()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Open database connection
	db, err := database.Open(cfg)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Perform migrations
	err = db.AutoMigrate(&models.User{}, &models.Dashboard{}, &models.Project{}, &models.Task{})
	if err != nil {
		log.Fatalf("couldn't auto migrate: %v", err)
	}

	// Initialise router
	r := mux.NewRouter()

	// Initialise middleware
	authMiddleware := token.NewAuthMiddleware(cfg.SigningKey)

	// Initialise services
	project.NewProjectService(r, *authMiddleware, repository.NewProjectRepository(db))
	task.NewTaskService(r, repository.NewTaskRepository(db), *authMiddleware)

	// Start the servers and listen
	log.Fatal(http.ListenAndServe(":8083", r))
}
