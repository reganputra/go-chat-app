package jwt

import (
	"context"
	"errors"
	"fmt"
	"go-chat-app/pkg/env"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.elastic.co/apm"
)

type ClaimToken struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	jwt.RegisteredClaims
}

var MapTokenTypes = map[string]time.Duration{
	"access":  time.Minute * 15,
	"refresh": time.Hour * 24,
}

func GenerateToken(ctx context.Context, username, fullName, tokenType string, now time.Time) (string, error) {

	span, _ := apm.StartSpan(ctx, "GenerateToken", "jwt")
	defer span.End()

	secret := []byte(env.GetEnv("APP_SECRET", ""))

	claims := ClaimToken{
		Username: username,
		FullName: fullName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    env.GetEnv("APP_NAME", ""),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(MapTokenTypes[tokenType])),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return tokenString, errors.New("failed to generate token")
	}
	return tokenString, nil
}

func ValidateToken(ctx context.Context, token string) (*ClaimToken, error) {

	span, _ := apm.StartSpan(ctx, "ValidateToken", "jwt")
	defer span.End()

	secret := []byte(env.GetEnv("APP_SECRET", ""))

	parsedToken, err := jwt.ParseWithClaims(token, &ClaimToken{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(*ClaimToken); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}
