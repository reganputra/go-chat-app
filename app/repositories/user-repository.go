package repositories

import (
	"context"
	"go-chat-app/app/models"
	"go-chat-app/pkg/database"
	"time"

	"go.elastic.co/apm"
)

func CreateUser(ctx context.Context, user *models.User) error {
	span, _ := apm.StartSpan(ctx, "CreateUser", "repository")
	defer span.End()

	return database.DB.Create(user).Error
}

func GetUserByUsername(ctx context.Context, username string) (models.User, error) {

	span, _ := apm.StartSpan(ctx, "GetUserByUsername", "repository")
	defer span.End()

	var user models.User
	return user, database.DB.Where("username = ?", username).Last(&user).Error
}

func CreateUserSession(ctx context.Context, session *models.UserSession) error {
	span, _ := apm.StartSpan(ctx, "CreateUserSession", "repository")
	defer span.End()

	return database.DB.Create(session).Error
}

func DeleteUserSession(ctx context.Context, token string) error {

	span, _ := apm.StartSpan(ctx, "DeleteUserSession", "repository")
	defer span.End()

	return database.DB.Exec("DELETE FROM user_sessions WHERE token = ?", token).Error
}

func GetUserSession(ctx context.Context, token string) (models.UserSession, error) {

	span, _ := apm.StartSpan(ctx, "GetUserSession", "repository")
	defer span.End()

	var session models.UserSession
	return session, database.DB.Where("token = ?", token).Last(&session).Error
}

func UpdateUserSessionTokens(ctx context.Context, accessToken, refreshToken string,
	tokenExpired, refreshTokenExpired time.Time, oldRefreshToken string) error {

	span, _ := apm.StartSpan(ctx, "UpdateUserSessionTokens", "repository")
	defer span.End()

	return database.DB.Exec(`UPDATE user_sessions 
        SET token = ?, refresh_token = ?, token_expired = ?, refresh_token_expired = ?, updated_at = ? 
        WHERE refresh_token = ?`,
		accessToken, refreshToken, tokenExpired, refreshTokenExpired, time.Now(), oldRefreshToken).Error
}

func GetUserSessionByRefreshToken(ctx context.Context, refreshToken string) (models.UserSession, error) {

	span, _ := apm.StartSpan(ctx, "GetUserSessionByRefreshToken", "repository")
	defer span.End()

	var session models.UserSession
	return session, database.DB.Where("refresh_token = ?", refreshToken).Last(&session).Error
}
