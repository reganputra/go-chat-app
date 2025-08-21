package repositories

import (
	"context"
	"go-chat-app/app/models"
	"go-chat-app/pkg/database"
)

func CreateUser(ctx context.Context, user *models.User) error {
	return database.DB.Create(user).Error
}

func GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User
	return user, database.DB.Where("username = ?", username).Last(&user).Error
}

func CreateUserSession(ctx context.Context, session *models.UserSession) error {
	return database.DB.Create(session).Error
}

func DeleteUserSession(ctx context.Context, token string) error {
	return database.DB.Exec("DELETE FROM user_sessions WHERE token = ?", token).Error
}

func GetUserSession(ctx context.Context, token string) (models.UserSession, error) {
	var session models.UserSession
	return session, database.DB.Where("token = ?", token).Last(&session).Error
}
