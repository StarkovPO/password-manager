package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"password-manager/internal/config"
	"password-manager/internal/service"
	"strings"
)

type Middleware struct {
	c config.Config
}

func NewMiddleware(c *config.Config) *Middleware {
	return &Middleware{c: *c}
}

func (m *Middleware) CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}
		tokenTrim := strings.TrimPrefix(token, "Bearer ")
		UID, err := m.parseToken(tokenTrim)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r.Header.Add("User-ID", UID)

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) parseToken(accessToken string) (string, error) {

	token, err := jwt.ParseWithClaims(accessToken, &service.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("invalid signing method")
		}
		return []byte(m.c.SecretValue), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*service.TokenClaims)
	if !ok {
		return "", errors.New("token claims are not of type")
	}
	return claims.UserID, nil
}
