package handler

import (
	"etalert-backend/service"

	"github.com/gofiber/fiber/v2"
)

type WeeklyReportListHandler struct {
	weeklyReportListsrv service.WeeklyReportListService
}

func NewWeeklyReportListHandler(weeklyReportListService service.WeeklyReportListService) *WeeklyReportListHandler {
	return &WeeklyReportListHandler{weeklyReportListsrv: weeklyReportListService}
}

func (h *WeeklyReportListHandler) GetWeeklyReportLists(c *fiber.Ctx) error {
	googleId := c.Params("googleId")

	weeklyReportList, err := h.weeklyReportListsrv.GetWeeklyReportLists(googleId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get schedule"})
	}
	if len(weeklyReportList) == 0 {
		return c.JSON([]interface{}{})
	}
	return c.JSON(weeklyReportList)
}
