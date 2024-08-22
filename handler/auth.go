package handler

import (
	"etalert-backend/service"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
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
    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing refresh token"})
    }

    // Expecting "Bearer <token>", so split by space and take the second part
    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid authorization header format"})
    }

    refreshToken := parts[1]

    refreshResponse, err := h.authsrv.RefreshToken(refreshToken)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
    }

    return c.Status(http.StatusOK).JSON(refreshResponse)
}


