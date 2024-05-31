package handler

import (
	"etalert-backend/service"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type userHandler struct {
	usersrv service.UserService
}

type createUserRequest struct {
	Name     string `json:"name"`
	Image    string `json:"image"`
	Email    string `json:"email"`
	GoogleId string `json:"googleId"`
}

type createUserResponse struct {
	Message string `json:"message"`
}

func NewUserHandler(userService service.UserService) *userHandler {
	return &userHandler{usersrv: userService}
}

func (h *userHandler) CreateUser(c *fiber.Ctx) error {
	var req createUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	user := &service.UserInput{
		Name:     req.Name,
		Image:    req.Image,
		Email:    req.Email,
		GoogleId: req.GoogleId,
	}

	err := h.usersrv.InsertUser(user)
	if err != nil {
		if err == service.ErrUserAlreadyExists {
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "User with the same GoogleId already exists"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert user"})
	}

	return c.Status(http.StatusCreated).JSON(createUserResponse{Message: "User created successfully"})
}

func (h *userHandler) GetUserInfo(c *fiber.Ctx) error {
	googleId := c.Params("googleId")

	user, err := h.usersrv.GetUserInfo(googleId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user"})
	}
	if user == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(http.StatusOK).JSON(user)
}
