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

// Get godoc
// @Summary Получить анализ по ID
// @Description Возвращает Analysis вместе со списком Findings
// @Tags Analysis
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID анализа"
// @Success 200 {object} dto.AnalysisResponse
// @Failure 401 {object} dto.ErrorResponse "Неавторизован"
// @Failure 404 {object} dto.ErrorResponse "Анализ не найден"
// @Router /analysis/{id} [get]

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
