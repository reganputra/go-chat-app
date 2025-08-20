package repositories

import (
	"context"
	"go-chat-app/app/models"
	"go-chat-app/pkg/database"
)

func CreateUser(ctx context.Context, user *models.User) error {
	return database.DB.Create(user).Error
}
