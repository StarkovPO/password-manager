package service

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const tokenTTL = 30 * time.Minute

type TokenClaims struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
}

func NewToken(ID string) *jwt.Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		ID,
	})

	return token
}
