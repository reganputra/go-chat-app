package database

import (
	"fmt"
	"go-chat-app/app/models"
	"go-chat-app/pkg/env"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDatabase() {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", env.GetEnv("DB_USER", ""), env.GetEnv("DB_PASSWORD", ""), env.GetEnv("DB_HOST", "127.0.0.1"), env.GetEnv("DB_PORT", "3306"), env.GetEnv("DB_NAME", ""))

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database! \n", err.Error())
		os.Exit(1)
	}

	err = DB.AutoMigrate(&models.User{}, &models.UserSession{})
	if err != nil {
		log.Fatal("Failed to migrate the Database! \n", err.Error())
	}

	DB.Logger = logger.Default.LogMode(logger.Info)
}

func SetupMongoDb() {
	uri := env.GetEnv("MONGODB_URI", "")
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	coll := client.Database("messaging_db").Collection("message_history")
	MongoDB = coll

	log.Println("successfully connected to mongoDB")
}
