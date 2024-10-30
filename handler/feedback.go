package handler

import (
	"etalert-backend/service"
	"etalert-backend/validators"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type FeedbackHandler struct {
	feedbacksrv service.FeedbackService
}

type createFeedbackRequest struct {
	GoogleId string   `json:"googleId" validate:"required"`
	Feedback string   `json:"feedback" validate:"required"`
}

type createFeedbackResponse struct {
	Message string `json:"message"`
}

func NewFeedbackHandler(feedbackService service.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{feedbacksrv: feedbackService}
}

func (h *FeedbackHandler) CreateFeedback(c *fiber.Ctx) error {
	var req createFeedbackRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := validators.ValidateStruct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	feedback := &service.FeedbackInput{
		GoogleId: req.GoogleId,
		Feedback: req.Feedback,
	}

	err := h.feedbacksrv.InsertFeedback(feedback)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert feedback"})
	}

	return c.Status(fiber.StatusCreated).JSON(createFeedbackResponse{Message: "Feedback created successfully"})
}