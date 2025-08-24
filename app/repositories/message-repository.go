package repositories

import (
	"context"
	"errors"
	"go-chat-app/app/models"
	"go-chat-app/pkg/database"

	"go.elastic.co/apm"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func InsertNewMessage(ctx context.Context, data models.MessagePayload) error {

	span, _ := apm.StartSpan(ctx, "InsertNewMessage", "repository")
	defer span.End()

	_, err := database.MongoDB.InsertOne(ctx, data)
	return err
}

func GetAllMessage(ctx context.Context) ([]models.MessagePayload, error) {

	span, _ := apm.StartSpan(ctx, "GetAllMessage", "repository")
	defer span.End()

	var msg []models.MessagePayload

	cursor, err := database.MongoDB.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.New("failed to get all messages")
	}

	for cursor.Next(ctx) {
		payload := models.MessagePayload{}
		err := cursor.Decode(&payload)
		if err != nil {
			return msg, errors.New("failed to decode message")
		}
		msg = append(msg, payload)
	}
	return msg, nil
}
