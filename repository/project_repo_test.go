package repository

//
//import (
//	"encoding/json"
//	"testing"
//
//	"github.com/stretchr/testify/require"
//	"github.com/stretchr/testify/suite"
//	"gorm.io/gorm"
//
//	"github.com/todanni/api/models"
//	"github.com/todanni/api/test"
//)
//
//type ProjectRepositoryTestSuite struct {
//	test.DbSuite
//}
//
//func (s *ProjectRepositoryTestSuite) SetupSuite() {
//	s.Db, s.CleanupFunc = test.SetupGormWithDocker()
//	s.Db.AutoMigrate(&models.Project{}, &models.User{})
//}
//
//func (s *ProjectRepositoryTestSuite) TearDownSuite() {
//	s.CleanupFunc()
//}
//
//func (s *ProjectRepositoryTestSuite) Test_CreateProject() {
//	repo := projectRepo{db: s.Db}
//	userRepo := userRepo{db: s.Db}
//
//	userOne, err := userRepo.CreateUser(models.User{
//		Email:      "userOne@mail.com",
//		FirstName:  "User",
//		LastName:   "One",
//		ProfilePic: "https://imgur.com/userone.png",
//	})
//	require.NoError(s.T(), err)
//	require.NotEqual(s.T(), 0, userOne.ID)
//
//	userTwo, err := userRepo.CreateUser(models.User{
//		Email:      "userTwo@mail.com",
//		FirstName:  "User",
//		LastName:   "Two",
//		ProfilePic: "https://imgur.com/usertwo.png",
//	})
//	require.NoError(s.T(), err)
//	require.Equal(s.T(), "userTwo@mail.com", userTwo.Email)
//
//	result, err := repo.CreateProject(models.Project{
//		Name:  "Name",
//		Owner: 1,
//		Members: []models.User{
//			{
//				Model: gorm.Model{
//					ID: userOne.ID,
//				},
//			},
//		},
//	})
//	require.NoError(s.T(), err)
//	require.NotNil(s.T(), result)
//
//	marshalled, err := json.Marshal(result)
//	require.NoError(s.T(), err)
//	require.NotNil(s.T(), marshalled)
//
//	err = repo.AddProjectMember(userTwo.ID, result.ID)
//	require.NoError(s.T(), err)
//
//	projectsForUser, err := repo.ListProjectsByUser(userOne.ID)
//	require.NoError(s.T(), err)
//	require.NotNil(s.T(), projectsForUser)
//
//	err = repo.RemoveProjectMember(userTwo.ID, result.ID)
//	require.NoError(s.T(), err)
//
//	projectsForUser, err = repo.ListProjectsByUser(userOne.ID)
//	require.NoError(s.T(), err)
//	require.NotNil(s.T(), projectsForUser)
//}
//
//func TestProjectRepositoryTestSuite(t *testing.T) {
//	suite.Run(t, new(ProjectRepositoryTestSuite))
//}
