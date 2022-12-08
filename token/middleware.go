package token

import (
	"context"
	"errors"
	"net/http"
	"strings"
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
			accessToken, err = m.checkCookieValue(r)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
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
	accessToken := &ToDanniToken{}
	err := accessToken.Parse(requestToken, m.signingKey)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (m *AuthMiddleware) checkCookieValue(r *http.Request) (*ToDanniToken, error) {
	accessTokenCookie, err := r.Cookie(AccessTokenCookieName)
	// If cookie is not present, check the authorization header
	if err != nil {
		return nil, errors.New("access token cookie wasn't set")
	}

	accessToken := &ToDanniToken{}
	err = accessToken.Parse(accessTokenCookie.Value, m.signingKey)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}
