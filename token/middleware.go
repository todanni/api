package token

import (
	"context"
	"errors"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	ErrorEmptyAuthHeader = errors.New("authorization header wasn't set")
	ErrorTokenNotPresent = errors.New("token not present")
)

type AuthMiddleware struct {
	signingKey string
}

func NewAuthMiddleware(signingKey string) *AuthMiddleware {
	return &AuthMiddleware{
		signingKey: signingKey,
	}
}

func (m *AuthMiddleware) JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := m.checkAuthHeader(r)
		if err != nil {
			log.Infof("couldn't find valid access token in cookie: %v", err)
			accessToken, err = m.checkCookieValue(r)
		}

		//accessToken, err := m.checkCookieValue(r)

		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), AccessTokenContextKey, accessToken)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) checkAuthHeader(r *http.Request) (*ToDanniToken, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, ErrorEmptyAuthHeader
	}

	parts := strings.Split(authHeader, "Bearer ")
	if len(parts) != 2 {
		return nil, ErrorTokenNotPresent
	}

	requestToken := parts[1]
	return m.parseToken(requestToken)
}

func (m *AuthMiddleware) checkCookieValue(r *http.Request) (*ToDanniToken, error) {
	accessTokenCookie, err := r.Cookie(AccessTokenCookieName)
	// If cookie is not present, check the authorization header

	if err != nil {
		log.Error(err)
		return nil, errors.New("access token cookie wasn't set")
	}
	return m.parseToken(accessTokenCookie.Value)
}

func (m *AuthMiddleware) parseToken(tokenString string) (*ToDanniToken, error) {
	accessToken := &ToDanniToken{}
	err := accessToken.Parse(tokenString, m.signingKey)
	return accessToken, err
}
