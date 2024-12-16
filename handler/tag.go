package handler

import (
	"etalert-backend/service"

	"github.com/gofiber/fiber/v2"
)

type TagHandler struct {
	tagsrv service.TagService
}

type createTagRequest struct {
	GoogleId string   `json:"googleId" validate:"required"`
	Name     string   `json:"name" validate:"required"`
	Routines []string `json:"routines"`
}

type createTagResponse struct {
	Message string `json:"message"`
}

func NewTagHandler(tagService service.TagService) *TagHandler {
	return &TagHandler{tagsrv: tagService}
}

func (h *TagHandler) CreateTag(c *fiber.Ctx) error {
	var req createTagRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	tag := &service.TagInput{
		GoogleId: req.GoogleId,
		Name:     req.Name,
		Routines: req.Routines,
	}

	err := h.tagsrv.InsertTag(tag)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert tag"})
	}

	return c.Status(fiber.StatusCreated).JSON(createTagResponse{Message: "Tag created successfully"})
}

func (h *TagHandler) GetAllTags(c *fiber.Ctx) error {
	googleId := c.Params("googleId")

	tags, err := h.tagsrv.GetAllTags(googleId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get tags"})
	}

	if len(tags) == 0 {
		return c.JSON([]interface{}{})
	}

	return c.Status(fiber.StatusOK).JSON(tags)
}

func (h *TagHandler) GetRoutinesByTagId(c *fiber.Ctx) error {
	id := c.Params("id")

	routines, err := h.tagsrv.GetRoutinesByTagId(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get routines"})
	}

	if len(routines) == 0 {
		return c.JSON([]interface{}{})
	}

	return c.Status(fiber.StatusOK).JSON(routines)
}

func (h *TagHandler) UpdateTag(c *fiber.Ctx) error {
	id := c.Params("id")

	var req createTagRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	tag := &service.TagUpdateInput{
		Name:     req.Name,
		Routines: req.Routines,
	}

	err := h.tagsrv.UpdateTag(id, tag)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update tag"})
	}

	return c.Status(fiber.StatusOK).JSON(createTagResponse{Message: "Tag updated successfully"})
}

func (h *TagHandler) DeleteTag(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.tagsrv.DeleteTag(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete tag"})
	}

	return c.Status(fiber.StatusOK).JSON(createTagResponse{Message: "Tag deleted successfully"})
}