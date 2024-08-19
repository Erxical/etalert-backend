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

	routineRepository := repository.NewRoutineRepositoryDB(client, "etalert", "routine")
	routineService := service.NewRoutineService(routineRepository)
	routineHandler := handler.NewRoutineHandler(routineService)

	authService := service.NewAuthService(userRepository)
	authHandler := handler.NewAuthHandler(authService)

	// initialize new instance of fiber
	server := fiber.New()

	server.Post("/login", authHandler.Login)
	server.Post("/users", userHandler.CreateUser)
	server.Patch("/users/:googleId", userHandler.UpdateUser)
	server.Get("/users/info/:googleId", userHandler.GetUserInfo)
	server.Post("/users/bedtimes", bedtimeHandler.CreateBedtime)
	server.Patch("/users/bedtimes/:googleId", bedtimeHandler.UpdateBedtime)
	server.Get("/users/bedtimes/info/:googleId", bedtimeHandler.GetBedtimeInfo)
	server.Post("/users/routines", routineHandler.CreateRoutine)
	server.Patch("/users/routines/:googleId", routineHandler.UpdateRoutine)
	server.Get("/users/routines/info/:googleId", routineHandler.GetRoutineInfo)

	// listen to port 3000
	log.Fatal(server.Listen("localhost:3000"))
}
