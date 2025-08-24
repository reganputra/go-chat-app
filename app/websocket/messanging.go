package websocket

import (
	"context"
	"fmt"
	"go-chat-app/app/models"
	"go-chat-app/app/repositories"
	"go-chat-app/pkg/env"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm"
)

func ServeWsMessage(app *fiber.App) {
	var clients = make(map[*websocket.Conn]bool)
	var broadcast = make(chan models.MessagePayload)

	app.Get("/message/v1/send", websocket.New(func(c *websocket.Conn) {
		defer func() {
			c.Close()
			delete(clients, c)
		}()
		clients[c] = true

		for {
			var msg models.MessagePayload
			err := c.ReadJSON(&msg)
			if err != nil {
				log.Printf("Error reading from client: %v", err)
				break
			}

			tx := apm.DefaultTracer.StartTransaction("Send Message", "websocket")
			ctx := apm.ContextWithTransaction(context.Background(), tx)

			msg.Date = time.Now()
			err = repositories.InsertNewMessage(ctx, msg)
			if err != nil {
				log.Printf("Error inserting message: %v", err)
				break
			}
			tx.End()
			broadcast <- msg
		}
	}))

	go func() {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error writing to client: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}()

	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", env.GetEnv("APP_HOST", "localhost"), env.GetEnv("APP_PORT_SOCKET", "8080"))))
}
