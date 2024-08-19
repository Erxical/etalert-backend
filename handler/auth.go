package handler

import (
	"etalert-backend/service"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type authHandler struct {
	authsrv service.AuthService
}

func NewAuthHandler(authService service.AuthService) *authHandler {
	return &authHandler{authsrv: authService}
}

func (h *authHandler) Login(c *fiber.Ctx) error {
	var req service.LoginInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	loginResponse, err := h.authsrv.Login(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to login"})
	}

	return c.Status(http.StatusOK).JSON(loginResponse)
}

func (h *authHandler) RefreshToken(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	refreshResponse, err := h.authsrv.RefreshToken(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
	}

	return c.Status(http.StatusOK).JSON(refreshResponse)
}
