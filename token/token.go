package token

import (
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	log "github.com/sirupsen/logrus"

	"github.com/todanni/api/models"
)

type ContextKey string

const (
	ToDanniTokenIssuer               = "todanni.com"
	AccessTokenCookieName            = "todanni-access-token"
	AccessTokenContextKey ContextKey = "accessToken"
)

var (
	ExpirationTime = 24 * time.Hour
)

type ToDanniToken struct {
	token jwt.Token
}

// NewAccessToken returns a ToDanni JWT issued at the current time
// with no claims yet set on it, other than issuer.
func NewAccessToken() *ToDanniToken {
	t, _ := jwt.NewBuilder().
		Issuer(ToDanniTokenIssuer).
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(ExpirationTime)).
		Build()
	return &ToDanniToken{token: t}
}

// SignToken returns the
func (t *ToDanniToken) SignToken(signingKey []byte) ([]byte, error) {
	signed, err := jwt.Sign(t.token, jwt.WithKey(jwa.HS256, signingKey))
	if err != nil {
		return nil, err
	}

	return signed, nil
}

func (t *ToDanniToken) Parse(signedToken, signingKey string) error {
	verifiedToken, err := jwt.Parse([]byte(signedToken), jwt.WithKey(jwa.HS256, []byte(signingKey)))
	if err != nil {
		log.Error(err)
		return err
	}

	t.token = verifiedToken
	return nil
}

func (t *ToDanniToken) SetUserID(id string) {
	t.setClaim("user_id", id)
}

func (t *ToDanniToken) GetUserID() string {
	userID, ok := t.token.Get("user_id")
	if !ok {
		return ""
	}
	return userID.(string)
}

func (t *ToDanniToken) SetDashboardPermissions(dashboards []models.Dashboard) *ToDanniToken {
	userDashboardIDs := make([]uuid.UUID, 0)

	for _, dashboard := range dashboards {
		userDashboardIDs = append(userDashboardIDs, dashboard.ID)
	}

	return t.setClaim("dashboards", userDashboardIDs)
}

func (t *ToDanniToken) HasDashboardPermission(dashboard uuid.UUID) bool {
	dashboardsPermissions, ok := t.token.Get("dashboards")
	if !ok {
		return false
	}

	dashboardsPermissionsArray := dashboardsPermissions.([]interface{})

	for _, value := range dashboardsPermissionsArray {
		if value.(string) == dashboard.String() {
			return true
		}
	}
	return false
}

func (t *ToDanniToken) SetProjectsPermissions(projects []models.Project) *ToDanniToken {
	userProjectIDs := make([]uint, 0)

	for _, project := range projects {
		userProjectIDs = append(userProjectIDs, project.ID)
	}

	return t.setClaim("projects", userProjectIDs)
}

func (t *ToDanniToken) HasProjectPermission(project uint) bool {
	projectPermissions, ok := t.token.Get("projects")
	if !ok {
		return false
	}

	projectsPermissionsArray := projectPermissions.([]interface{})

	for _, permission := range projectsPermissionsArray {
		projectID := permission.(float64)

		if uint(projectID) == project {
			return true
		}
	}
	return false
}

func (t *ToDanniToken) setClaim(name string, value interface{}) *ToDanniToken {
	_ = t.token.Set(name, value)
	return t
}
