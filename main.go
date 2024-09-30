package main

import (
	"context"
	"etalert-backend/handler"
	"etalert-backend/repository"
	"etalert-backend/service"
	"etalert-backend/middlewares"
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
	
	authService := service.NewAuthService(userRepository)
	authHandler := handler.NewAuthHandler(authService)

	bedtimeRepository := repository.NewBedtimeRepositoryDB(client, "etalert", "bedtime")
	bedtimeService := service.NewBedtimeService(bedtimeRepository)
	bedtimeHandler := handler.NewBedtimeHandler(bedtimeService)

	routineRepository := repository.NewRoutineRepositoryDB(client, "etalert", "routine")
	routineService := service.NewRoutineService(routineRepository)
	routineHandler := handler.NewRoutineHandler(routineService)

	scheduleRepository := repository.NewScheduleRepositoryDB(client, "etalert", "schedule")
	scheduleService := service.NewScheduleService(scheduleRepository, routineRepository, bedtimeRepository)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)

	scheduleService.StartCronJob()

	// initialize new instance of fiber
	server := fiber.New()

	server.Post("/login", authHandler.Login)
	server.Post("/refresh-token", authHandler.RefreshToken)
    server.Post("/create-user", userHandler.CreateUser)

	protected := server.Group("/users", middlewares.ValidateSession(authService))

    // Protected routes
	//User routes
    protected.Patch("/:googleId", userHandler.UpdateUser)
    protected.Get("/info/:googleId", userHandler.GetUserInfo)

	//Bedtime routes
    protected.Post("/bedtimes", bedtimeHandler.CreateBedtime)
    protected.Patch("/bedtimes/:googleId", bedtimeHandler.UpdateBedtime)
    protected.Get("/bedtimes/info/:googleId", bedtimeHandler.GetBedtimeInfo)

	//Routine routes
    protected.Post("/routines", routineHandler.CreateRoutine)
    protected.Patch("/routines/edit/:googleId", routineHandler.UpdateRoutine)
    protected.Get("/routines/:googleId", routineHandler.GetAllRoutines)

	//Schedule routes
	protected.Post("/schedules", scheduleHandler.CreateSchedule)
	protected.Get("/schedules/:googleId/:date?", scheduleHandler.GetAllSchedules)
	protected.Get("/schedules/:id", scheduleHandler.GetScheduleById)
	protected.Patch("/schedules/:id", scheduleHandler.UpdateSchedule)
	protected.Delete("/schedules/:groupId", scheduleHandler.DeleteSchedule)

	// listen to port 3000
	log.Fatal(server.Listen(":3000"))
}
