package jwt

import (
	"context"
	"errors"
	"go-chat-app/pkg/env"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ClaimToken struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	jwt.RegisteredClaims
}

var mapTokenTypes = map[string]time.Duration{
	"access":  time.Minute * 15,
	"refresh": time.Hour * 24,
}

func GenerateToken(ctx context.Context, username, fullName, tokenType string) (string, error) {
	secret := []byte(env.GetEnv("APP_SECRET", ""))

	claims := ClaimToken{
		Username: username,
		FullName: fullName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    env.GetEnv("APP_NAME", ""),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(mapTokenTypes[tokenType])),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return tokenString, errors.New("failed to generate token")
	}
	return tokenString, nil
}
