package analysis

import (
	"mws-ai/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AnalysisHandler struct {
	service *services.AnalysisService
}

func NewAnalysisHandler(service *services.AnalysisService) *AnalysisHandler {
	return &AnalysisHandler{
		service: service,
	}
}

func (h *AnalysisHandler) Get() fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fiber.ErrBadRequest
		}

		analysis, err := h.service.GetByID(uint(id))
		if err != nil {
			return fiber.ErrNotFound
		}

		return c.JSON(analysis)
	}
}
