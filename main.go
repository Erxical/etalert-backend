package main

import (
	"context"
	"etalert-backend/handler"
	"etalert-backend/middlewares"
	"etalert-backend/repository"
	"etalert-backend/service"
	etalert_websocket "etalert-backend/websocket"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	routineLogRepository := repository.NewRoutineLogRepositoryDB(client, "etalert", "routineLog")
	routineLogService := service.NewRoutineLogService(routineLogRepository)
	routineLogHandler := handler.NewRoutineLogHandler(routineLogService)

	routineRepository := repository.NewRoutineRepositoryDB(client, "etalert", "routine")
	tagRepository := repository.NewTagRepositoryDB(client, "etalert", "tag")
	
	routineService := service.NewRoutineService(routineRepository, tagRepository)
	routineHandler := handler.NewRoutineHandler(routineService)

	tagService := service.NewTagService(tagRepository, routineRepository)
	tagHandler := handler.NewTagHandler(tagService)

	weeklyReportListRepository := repository.NewWeeklyReportListRepositoryDB(client, "etalert", "weeklyReportList")
	weeklyReportListService := service.NewWeeklyReportListService(weeklyReportListRepository)
	weeklyReportListHandler := handler.NewWeeklyReportListHandler(weeklyReportListService)

	weeklyReportRepository := repository.NewWeeklyReportRepositoryDB(client, "etalert", "weeklyReport")
	weeklyReportService := service.NewWeeklyReportService(weeklyReportRepository, userRepository, routineRepository, weeklyReportListRepository, routineLogRepository, tagRepository)
	weeklyReportHandler := handler.NewWeeklyReportHandler(weeklyReportService)

	scheduleLogRepository := repository.NewScheduleLogRepositoryDB(client, "etalert", "scheduleLog")

	scheduleRepository := repository.NewScheduleRepositoryDB(client, "etalert", "schedule")
	scheduleService := service.NewScheduleService(scheduleRepository, scheduleLogRepository, routineRepository, bedtimeRepository, tagRepository)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)

	feedbackRepository := repository.NewFeedbackRepositoryDB(client, "etalert", "feedback")
	feedbackService := service.NewFeedbackService(feedbackRepository)
	feedbackHandler := handler.NewFeedbackHandler(feedbackService)

	scheduleService.StartCronJob()
	weeklyReportService.StartCronJob()

	// initialize new instance of fiber
	server := fiber.New()

	server.Get("/ws", websocket.New(func(c *websocket.Conn) {
		etalert_websocket.HandleConnections(c)
	}))

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
	protected.Get("/routines/:googleId", routineHandler.GetAllRoutines)
	protected.Patch("/routines/edit/:id", routineHandler.UpdateRoutine)
	protected.Delete("/routines/:id", routineHandler.DeleteRoutine)

	//RoutineLog routes
	protected.Post("/routine-logs", routineLogHandler.InsertRoutineLog)
	protected.Get("/routine-logs/:googleId/:date?", routineLogHandler.GetRoutineLogs)
	protected.Delete("/routine-logs/:id", routineLogHandler.DeleteRoutineLog)

	//Tag routes
	protected.Post("/tags", tagHandler.CreateTag)
	protected.Get("/tags/:googleId", tagHandler.GetAllTags)
	protected.Get("/tags/routines/:id", tagHandler.GetRoutinesByTagId)
	protected.Patch("/tags/:id", tagHandler.UpdateTag)
	protected.Delete("/tags/:id", tagHandler.DeleteTag)

	//WeeklyReport routes
	protected.Get("/weekly-reports/:googleId/:date", weeklyReportHandler.GetWeeklyReports)

	//WeeklyReportList routes
	protected.Get("/weekly-report-lists/:googleId", weeklyReportListHandler.GetWeeklyReportLists)

	//Schedule routes
	protected.Post("/schedules", scheduleHandler.CreateSchedule)
	protected.Get("/schedules/all/:googleId/:date?", scheduleHandler.GetAllSchedules)
	protected.Get("/schedules/:id", scheduleHandler.GetScheduleById)
	protected.Get("/schedules/group/:groupId", scheduleHandler.GetSchedulesByGroupId)
	protected.Get("/schedules/recurrence/:recurrenceId/:date?", scheduleHandler.GetSchedulesIdByRecurrenceId)
	protected.Patch("/schedules/:id", scheduleHandler.UpdateSchedule)
	protected.Patch(("/schedules/recurrence/:recurrenceId/:date?"), scheduleHandler.UpdateScheduleByRecurrenceId)
	protected.Delete("/schedules/:groupId", scheduleHandler.DeleteSchedule)
	protected.Delete("/schedules/recurrence/:recurrenceId/:date?", scheduleHandler.DeleteScheduleByRecurrenceId)

	//Feedback routes
	protected.Post("/create-feedbacks", feedbackHandler.CreateFeedback)

	// listen to port 3000
	log.Fatal(server.Listen(":3000"))
}
