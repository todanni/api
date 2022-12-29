package repository

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/todanni/api/models"
	"github.com/todanni/api/test"
)

type UserRepositoryTestSuite struct {
	test.DbSuite
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	s.Db, s.CleanupFunc = test.SetupGormWithDocker()
	s.Db.AutoMigrate(&models.User{})
}

func (s *UserRepositoryTestSuite) TearDownSuite() {
	s.CleanupFunc()
}

func (s *UserRepositoryTestSuite) Test_CreateUser() {
	repo := userRepo{db: s.Db}
	email := "user@email.com"

	result, err := repo.CreateUser(models.User{
		Email:      email,
		ProfilePic: "https://imgur.com/me.png",
	})
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)

	marshalled, err := json.Marshal(result)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), marshalled)

	user, err := repo.GetUserByEmail(email)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), user)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
