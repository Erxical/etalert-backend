package main

import (
	"context"
	"etalert-backend/handler"
	"etalert-backend/repository"
	"etalert-backend/service"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func main() {
	// initialize godotenv to read all .env files
	godotenv.Load()

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	userRepository := repository.NewUserRepositoryDB(client, "etalert", "user")
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	bedtimeRepository := repository.NewBedtimeRepositoryDB(client, "etalert", "bedtime")
	bedtimeService := service.NewBedtimeService(bedtimeRepository)
	bedtimeHandler := handler.NewBedtimeHandler(bedtimeService)

	// initialize new instance of fiber
	server := fiber.New()

	server.Post("/users", userHandler.CreateUser)
	server.Get("/users/:googleId", userHandler.GetUserInfo)
	server.Post("/users/bedtimes", bedtimeHandler.CreateBedtime)
	server.Get("/users/bedtimes/info/:googleId", bedtimeHandler.GetBedtimeInfo)
	server.Post("/users/bedtimes/:googleId", bedtimeHandler.UpdateBedtime)

	// listen to port 8000
	log.Fatal(server.Listen("localhost:8000"))
}
