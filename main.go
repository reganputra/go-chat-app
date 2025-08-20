package main

import (
	"fmt"
	"go-chat-app/bootstrap"
	"go-chat-app/pkg/env"
	"log"
)

func main() {

	app := bootstrap.NewApplication()
	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", env.GetEnv("APP_HOST", "localhost"), env.GetEnv("APP_PORT", "4000"))))

}
