package handler

import (
	"etalert-backend/service"
	"etalert-backend/validators"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type RoutineLogHandler struct {
	routineLogsrv service.RoutineLogService
}

type createRoutineLogRequest struct {
	RoutineId     string `json:"routineId" validate:"required"`
	GoogleId      string `json:"googleId" validate:"required"`
	Date          string `json:"date" validate:"required"`
	StartTime     string `json:"startTime" validate:"required"`
	EndTime       string `json:"endTime" validate:"required"`
	ActualEndTime string `json:"actualEndTime" validate:"required"`
	Skewness      int    `json:"skewness" validate:"required"`
}

type createRoutineLogResponse struct {
	Message string `json:"message"`
}

func NewRoutineLogHandler(routineLogService service.RoutineLogService) *RoutineLogHandler {
	return &RoutineLogHandler{routineLogsrv: routineLogService}
}

func (r *RoutineLogHandler) InsertRoutineLog(c *fiber.Ctx) error {
	var req createRoutineLogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := validators.ValidateStruct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	routineLog := &service.RoutineLogInput{
		RoutineId:     req.RoutineId,
		GoogleId:      req.GoogleId,
		Date:          req.Date,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		ActualEndTime: req.ActualEndTime,
		Skewness:      req.Skewness,
	}

	err := r.routineLogsrv.InsertRoutineLog(routineLog)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert routine log"})
	}

	return c.Status(fiber.StatusCreated).JSON(createRoutineLogResponse{Message: "Routine log created successfully"})
}

func (r *RoutineLogHandler) GetRoutineLogs(c *fiber.Ctx) error {
	googleId := c.Params("googleId")
	date := c.Params("date")

	routineLogs, err := r.routineLogsrv.GetRoutineLogs(googleId, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get routine logs"})
	}
	if len(routineLogs) == 0 {
		return c.JSON([]interface{}{})
	}

	return c.JSON(routineLogs)
}

func (r *RoutineLogHandler) DeleteRoutineLog(c *fiber.Ctx) error {
	id := c.Params("id")

	err := r.routineLogsrv.DeleteRoutineLog(id)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete routine log"})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}