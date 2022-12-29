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
		accessTokenString, err := m.checkCookieValue(r)
		if accessTokenString == "" {
			accessTokenString, err = m.checkAuthHeader(r)
		}

		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		log.Infof("Acess token: %s", accessTokenString)
		accessToken, err := m.parseToken(accessTokenString)
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

func (m *AuthMiddleware) checkAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrorEmptyAuthHeader
	}

	parts := strings.Split(authHeader, "Bearer ")
	if len(parts) != 2 {
		return "", ErrorTokenNotPresent
	}

	log.Info("Token Bearer found in request")
	return parts[1], nil
}

func (m *AuthMiddleware) checkCookieValue(r *http.Request) (string, error) {
	accessTokenCookie, err := r.Cookie(AccessTokenCookieName)
	if err != nil {
		return "", err
	}
	return accessTokenCookie.Value, err
}

func (m *AuthMiddleware) parseToken(tokenString string) (*ToDanniToken, error) {
	accessToken := &ToDanniToken{}
	err := accessToken.Parse(tokenString, m.signingKey)
	return accessToken, err
}
