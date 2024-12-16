package handler

import (
	"etalert-backend/service"

	"github.com/gofiber/fiber/v2"
)

type WeeklyReportHandler struct {
	weeklyReportsrv service.WeeklyReportService
}

func NewWeeklyReportHandler(weeklyReportService service.WeeklyReportService) *WeeklyReportHandler {
	return &WeeklyReportHandler{weeklyReportsrv: weeklyReportService}
}

func (h *WeeklyReportHandler) GetWeeklyReports(c *fiber.Ctx) error {
	googleId := c.Params("googleId")
	date := c.Params("date")

	weeklyReports, err := h.weeklyReportsrv.GetWeeklyReports(googleId, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get schedule"})
	}
	if len(weeklyReports) == 0 {
		return c.JSON([]interface{}{})
	}
	return c.JSON(weeklyReports)
}
