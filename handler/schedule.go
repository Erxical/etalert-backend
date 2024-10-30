package handler

import (
	"etalert-backend/service"
	"etalert-backend/validators"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ScheduleHandler struct {
	schedulesrv service.ScheduleService
}

type createScheduleRequest struct {
	GoogleId        string  `json:"googleId" validate:"required"`
	Name            string  `json:"name" validate:"required"`
	Date            string  `json:"date" validate:"required"`
	StartTime       string  `json:"startTime" validate:"required"`
	EndTime         string  `json:"endTime"`
	IsHaveEndTime   bool    `json:"isHaveEndTime"`
	OriName         string  `json:"oriName"`
	OriLatitude     float64 `json:"oriLatitude"`
	OriLongitude    float64 `json:"oriLongitude"`
	DestName        string  `json:"destName"`
	DestLatitude    float64 `json:"destLatitude"`
	DestLongitude   float64 `json:"destLongitude"`
	Priority        int     `json:"priority"`
	IsHaveLocation  bool    `json:"isHaveLocation"`
	IsFirstSchedule bool    `json:"isFirstSchedule"`

	Recurrence      string  `json:"recurrence"`
	RecurrenceId    int     `json:"recurrenceId"`
}

type updateScheduleRequest struct {
	Name          string `json:"name" validate:"required"`
	Date          string `json:"date" validate:"required"`
	StartTime     string `json:"startTime" validate:"required"`
	EndTime       string `json:"endTime"`
	IsHaveEndTime bool   `json:"isHaveEndTime"`
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
		GoogleId:        req.GoogleId,
		Name:            req.Name,
		Date:            req.Date,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		IsHaveEndTime:   req.IsHaveEndTime,
		OriName:         req.OriName,
		OriLatitude:     req.OriLatitude,
		OriLongitude:    req.OriLongitude,
		DestName:        req.DestName,
		DestLatitude:    req.DestLatitude,
		DestLongitude:   req.DestLongitude,
		Priority:        req.Priority,
		IsHaveLocation:  req.IsHaveLocation,
		IsFirstSchedule: req.IsFirstSchedule,

		Recurrence:      req.Recurrence,
	}

	if schedule.Recurrence != "none" {
		str, err := h.schedulesrv.InsertRecurrenceSchedule(schedule)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert schedule"})
		}
		if str != "" {
			return c.Status(fiber.StatusCreated).JSON(createScheduleResponse{Message: "Schedule created successfully with warning " + str})
		}
	
		return c.Status(fiber.StatusCreated).JSON(createScheduleResponse{Message: "Schedule created successfully"})
	}

	str, err := h.schedulesrv.InsertSchedule(schedule)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert schedule"})
	}

	if str != "" {
		return c.Status(fiber.StatusCreated).JSON(createScheduleResponse{Message: "Schedule created successfully with warning " + str})
	}

	return c.Status(fiber.StatusCreated).JSON(createScheduleResponse{Message: "Schedule created successfully"})
}

func (h *ScheduleHandler) GetAllSchedules(c *fiber.Ctx) error {
	googleId := c.Params("googleId")
	date := c.Params("date", "")

	schedules, err := h.schedulesrv.GetAllSchedules(googleId, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get schedule"})
	}
	if len(schedules) == 0 {
		return c.JSON([]interface{}{})
	}
	return c.JSON(schedules)
}

func (h *ScheduleHandler) GetScheduleById(c *fiber.Ctx) error {
	id := c.Params("id")

	schedule, err := h.schedulesrv.GetScheduleById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get schedule"})
	}
	if schedule == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Schedule not found"})
	}
	return c.JSON(schedule)
}

func (h *ScheduleHandler) UpdateSchedule(c *fiber.Ctx) error {
	id := c.Params("id")
	var req updateScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := validators.ValidateStruct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	schedule := &service.ScheduleUpdateInput{
		Name:          req.Name,
		Date:          req.Date,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		IsHaveEndTime: req.IsHaveEndTime,
	}

	err := h.schedulesrv.UpdateSchedule(id, schedule)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update schedule"})
	}

	return c.JSON(createScheduleResponse{Message: "Schedule updated successfully"})
}

func (h *ScheduleHandler) UpdateScheduleByRecurrenceId(c *fiber.Ctx) error {
	recurrenceId := c.Params("recurrenceId")
	date := c.Params("date", "")
	var req updateScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := validators.ValidateStruct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	schedule := &service.ScheduleUpdateInput{
		Name:          req.Name,
		Date:          req.Date,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		IsHaveEndTime: req.IsHaveEndTime,
	}

	err := h.schedulesrv.UpdateScheduleByRecurrenceId(recurrenceId, schedule, date)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update schedule"})
	}

	return c.JSON(createScheduleResponse{Message: "Schedule updated successfully"})
}

func (h *ScheduleHandler) DeleteSchedule(c *fiber.Ctx) error {
	groupId := c.Params("groupId")

	err := h.schedulesrv.DeleteSchedule(groupId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete schedule"})
	}

	return c.JSON(createScheduleResponse{Message: "Schedule deleted successfully"})
}

func (h *ScheduleHandler) DeleteScheduleByRecurrenceId(c *fiber.Ctx) error {
	recurrenceId := c.Params("recurrenceId")
	date := c.Params("date", "")

	err := h.schedulesrv.DeleteScheduleByRecurrenceId(recurrenceId, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete schedule"})
	}

	return c.JSON(createScheduleResponse{Message: "Schedule deleted successfully"})
}