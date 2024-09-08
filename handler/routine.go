package handler

import (
	"etalert-backend/service"
	"etalert-backend/validators"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type RoutineHandler struct {
	routinesrv service.RoutineService
}

type createRoutineRequest struct {
	GoogleId string `json:"googleId" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Duration int    `json:"duration" validate:"required"`
	Order    int    `json:"order" validate:"required"`
}

type createRoutineResponse struct {
	Message string `json:"message"`
}

func NewRoutineHandler(routineService service.RoutineService) *RoutineHandler {
	return &RoutineHandler{routinesrv: routineService}
}

func (h *RoutineHandler) CreateRoutine(c *fiber.Ctx) error {
	var req createRoutineRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := validators.ValidateStruct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	routine := &service.RoutineInput{
		GoogleId: req.GoogleId,
		Name:     req.Name,
		Duration: req.Duration,
		Order:    req.Order,
	}

	err := h.routinesrv.InsertRoutine(routine)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert routine"})
	}

	return c.Status(fiber.StatusCreated).JSON(createRoutineResponse{Message: "Routine created successfully"})
}

func (h *RoutineHandler) GetAllRoutines(c *fiber.Ctx) error {
	googleId := c.Params("googleId")

	routine, err := h.routinesrv.GetAllRoutines(googleId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get routine"})
	}
	if routine == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Routine not found"})
	}
	return c.JSON(routine)
}

func (h *RoutineHandler) UpdateRoutine(c *fiber.Ctx) error {
	googleId := c.Params("googleId")
	var req createRoutineRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := validators.ValidateStruct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	routine := &service.RoutineResponse{
		Name:     req.Name,
		Duration: req.Duration,
		Order:    req.Order,
	}

	err := h.routinesrv.UpdateRoutine(googleId, routine)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update routine"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Routine updated successfully"})
}
