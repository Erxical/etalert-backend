package handler

import (
	"etalert-backend/service"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

func (h *authHandler) CheckUserSession(c *fiber.Ctx) error {
    // Extract the access token from the Authorization header
    accessToken := c.Get("Authorization")
    if accessToken == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing access token"})
    }

    // Validate the access token
    claims, err := h.authsrv.ValidateAccessToken(accessToken)
    if err != nil {
        if err == jwt.ErrTokenExpired {
            // Access token expired, check if refresh token is provided
            refreshToken := c.Get("Authorization-Refresh")
            if refreshToken == "" {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Session expired and no refresh token provided"})
            }

            // Validate and refresh the access token using the refresh token
            refreshResponse, err := h.authsrv.RefreshToken(refreshToken)
            if err != nil {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
            }

            // Return the new access token to the client
            return c.Status(http.StatusOK).JSON(refreshResponse)
        }

        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid access token"})
    }

    // Token is valid, proceed with the request
    userId := claims["googleId"].(string)
    return c.Status(http.StatusOK).JSON(fiber.Map{"message": "User session is valid", "userId": userId})
}



