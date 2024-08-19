package handler

import (
	"etalert-backend/service"
	"etalert-backend/validators"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type bedtimeHandler struct {
	bedtimesrv service.BedtimeService
}

type createBedtimeRequest struct {
	GoogleId  string `json:"googleId" validate:"required"`
	SleepTime string `json:"sleepTime" validate:"required"`
	WakeTime  string `json:"wakeTime" validate:"required"`
}

type updateBedtimeRequest struct {
	SleepTime string `json:"sleepTime" validate:"required"`
	WakeTime  string `json:"wakeTime" validate:"required"`
}

type createBedtimeResponse struct {
	Message string `json:"message"`
}

func NewBedtimeHandler(bedtimeService service.BedtimeService) *bedtimeHandler {
	return &bedtimeHandler{bedtimesrv: bedtimeService}
}

func (h *bedtimeHandler) CreateBedtime(c *fiber.Ctx) error {
	var req createBedtimeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := validators.ValidateStruct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	bedtime := &service.BedtimeInput{
		GoogleId:  req.GoogleId,
		SleepTime: req.SleepTime,
		WakeTime:  req.WakeTime,
	}

	err := h.bedtimesrv.InsertBedtime(bedtime)
	if err != nil {
		if err == service.ErrBedtimeAlreadyExists {
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "Bedtime with the same GoogleId already exists"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert bedtime"})
	}

	return c.Status(http.StatusCreated).JSON(createBedtimeResponse{Message: "Bedtime created successfully"})
}

func (h *bedtimeHandler) GetBedtimeInfo(c *fiber.Ctx) error {
	googleId := c.Params("googleId")

	bedtime, err := h.bedtimesrv.GetBedtimeInfo(googleId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user"})
	}
	if bedtime == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Bedtime not found"})
	}

	return c.Status(http.StatusOK).JSON(bedtime)
}

func (h *bedtimeHandler) UpdateBedtime(c *fiber.Ctx) error {
	googleId := c.Params("googleId")

	var req updateBedtimeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := validators.ValidateStruct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	bedtime := &service.BedtimeResponse{
		SleepTime: req.SleepTime,
		WakeTime:  req.WakeTime,
	}

	err := h.bedtimesrv.UpdateBedtime(googleId, bedtime)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update bedtime"})
	}

	return c.Status(http.StatusOK).JSON(createBedtimeResponse{Message: "Bedtime updated successfully"})
}
