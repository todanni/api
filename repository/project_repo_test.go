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

type ProjectRepositoryTestSuite struct {
	test.DbSuite
}

func (s *ProjectRepositoryTestSuite) SetupSuite() {
	s.Db, s.CleanupFunc = test.SetupGormWithDocker()
	s.Db.AutoMigrate(&models.Project{}, &models.User{})
}

func (s *ProjectRepositoryTestSuite) TearDownSuite() {
	s.CleanupFunc()
}

func (s *ProjectRepositoryTestSuite) Test_CreateProject() {
	repo := projectRepo{db: s.Db}

	result, err := repo.CreateProject(models.Project{
		Name:  "Name",
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
	require.NotNil(s.T(), result)

	marshalled, err := json.Marshal(result)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), marshalled)
}

func TestProjectRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ProjectRepositoryTestSuite))
}
