package handler

import (
	"etalert-backend/service"
	"etalert-backend/validators"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ScheduleHandler struct {
	schedulesrv service.ScheduleService
}

type createScheduleRequest struct {
	GoogleId string `json:"googleId" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Date     string `json:"date" validate:"required"`
	StartTime string `json:"startTime" validate:"required"`
	EndTime string `json:"endTime"`
	IsHaveEndTime bool `json:"isHaveEndTime" validate:"required"`
	OriLatitude float64 `json:"oriLatitude"`
	OriLongitude float64 `json:"oriLongitude"`
	DestLatitude float64 `json:"destLatitude"`
	DestLongitude float64 `json:"destLongitude"`
	IsHaveLocation bool `json:"isHaveLocation" validate:"required"`
	IsFirstSchedule bool `json:"isFirstSchedule" validate:"required"`
	DepartTime string `json:"departTime"`
}

type createScheduleResponse struct {
	Message string `json:"message"`
}

func NewScheduleHandler(scheduleService service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{schedulesrv: scheduleService}
}

func (h *ScheduleHandler) CreateSchedule(c *fiber.Ctx) error {
	var req createScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := validators.ValidateStruct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	schedule := &service.ScheduleInput{
		GoogleId: req.GoogleId,
		Name:     req.Name,
		Date:     req.Date,
		StartTime: req.StartTime,
		EndTime: req.EndTime,
		IsHaveEndTime: req.IsHaveEndTime,
		OriLatitude: req.OriLatitude,
		OriLongitude: req.OriLongitude,
		DestLatitude: req.DestLatitude,
		DestLongitude: req.DestLongitude,
		IsHaveLocation: req.IsHaveLocation,
		IsFirstSchedule: req.IsFirstSchedule,
		DepartTime: req.DepartTime,
	}

	err := h.schedulesrv.InsertSchedule(schedule)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert schedule"})
	}

	return c.Status(fiber.StatusCreated).JSON(createScheduleResponse{Message: "Schedule created successfully"})
}