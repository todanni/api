package token

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/todanni/api/models"
)

const (
	signingKey = "examplesigningkey"
)

func TestToken_EndToEnd_NoDashboardsAndProjects(t *testing.T) {
	// Issue a token
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
	require.NoError(t, err)
	require.NotNil(t, signedToken)

	parsedToken := &ToDanniToken{}
	err = parsedToken.Parse(string(signedToken), signingKey)
	require.NoError(t, err)

	userID := parsedToken.GetUserID()
	require.Equal(t, id, userID)

	dashboardIDIsAllowed := parsedToken.HasDashboardPermission(uuid.New())
	require.Equal(t, false, dashboardIDIsAllowed)

	projectIDIsAllowed := parsedToken.HasProjectPermission(1)
	require.Equal(t, false, projectIDIsAllowed)

}

func TestToken_EndToEnd_OneDashboardAndProject(t *testing.T) {
	// Issue a token
	accessToken := NewAccessToken()
	id := uuid.New().String()
	id = id[:8]
	accessToken.SetUserID(id)

	// Set a token with no projects
	projects := []models.Project{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Name:  "Project",
			Owner: id,
			Members: []models.User{
				{
					ID: id,
				},
			},
		},
	}
	accessToken.SetProjectsPermissions(projects)
	dashboardID := uuid.New()

	// Set a token with one dashboard
	dashboards := []models.Dashboard{
		{
			ID: dashboardID,
		},
	}
	accessToken.SetDashboardPermissions(dashboards)

	signedToken, err := accessToken.SignToken([]byte(signingKey))
	require.NoError(t, err)
	require.NotNil(t, signedToken)

	parsedToken := &ToDanniToken{}
	err = parsedToken.Parse(string(signedToken), signingKey)
	require.NoError(t, err)

	userID := parsedToken.GetUserID()
	require.Equal(t, id, userID)

	dashboardIDIsAllowed := parsedToken.HasDashboardPermission(dashboardID)
	require.Equal(t, true, dashboardIDIsAllowed)

	projectIDIsAllowed := parsedToken.HasProjectPermission(1)
	require.Equal(t, true, projectIDIsAllowed)

}
