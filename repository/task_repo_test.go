package repository

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/todanni/api/models"
	"github.com/todanni/api/test"
)

type TaskRepositoryTestSuite struct {
	test.DbSuite
}

func (s *TaskRepositoryTestSuite) SetupSuite() {
	s.Db, s.CleanupFunc = test.SetupGormWithDocker()
	s.Db.AutoMigrate(&models.Task{}, &models.Project{}, &models.User{})
}

func (s *TaskRepositoryTestSuite) TearDownSuite() {
	s.CleanupFunc()
}

func (s *TaskRepositoryTestSuite) Test_CreateTask() {
	repo := taskRepo{db: s.Db}

	prjRepo := projectRepo{db: s.Db}

	projectOne, err := prjRepo.CreateProject(models.Project{
		Name:  "Project One",
		Owner: 1,
		Members: []models.User{
			{
				Model: gorm.Model{
					ID: 1,
				},
			},
		},
	})

	require.NoError(s.T(), err)
	require.NotNil(s.T(), projectOne)

	projectTwo, err := prjRepo.CreateProject(models.Project{
		Name:  "Project Two",
		Owner: 1,
		Members: []models.User{
			{
				Model: gorm.Model{
					ID: 1,
				},
			},
			{
				Model: gorm.Model{
					ID: 2,
				},
			},
		},
	})

	require.NoError(s.T(), err)
	require.NotNil(s.T(), projectTwo)

	projects, err := prjRepo.ListProjectsByUser(1)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), projects)

	taskOne, err := repo.CreateTask(models.Task{
		Title:     "Task Title",
		ProjectID: projectOne.ID,
		CreatedBy: 1,
	})
	require.NoError(s.T(), err)
	require.NotNil(s.T(), taskOne)

	taskTwo, err := repo.CreateTask(models.Task{
		Title:      "Task Title",
		ProjectID:  projectOne.ID,
		CreatedBy:  2,
		AssignedTo: 1,
	})
	require.NoError(s.T(), err)
	require.NotNil(s.T(), taskTwo)

	taskThree, err := repo.CreateTask(models.Task{
		Title:      "Task Title",
		ProjectID:  projectOne.ID,
		CreatedBy:  2,
		AssignedTo: 2,
	})
	require.NoError(s.T(), err)
	require.NotNil(s.T(), taskThree)

	tasks, err := repo.ListTasksByUser(1)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), tasks)

	marshalled, err := json.Marshal(tasks)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), marshalled)
}

func TestTaskRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TaskRepositoryTestSuite))
}
