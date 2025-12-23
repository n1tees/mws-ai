package analysis

import (
	"strconv"

	"mws-ai/internal/services"
	"mws-ai/pkg/logger"

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

		log := logger.Log.With().
			Str("handler", "analysis.get").
			Str("path", c.Path()).
			Logger()

		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Warn().
				Str("id_param", idStr).
				Msg("invalid analysis id")

			return fiber.ErrBadRequest
		}

		analysis, findings, err := h.service.GetDetails(uint(id))
		if err != nil || analysis == nil {
			log.Info().
				Int("analysis_id", id).
				Msg("analysis not found")

			return fiber.ErrNotFound
		}

		return c.JSON(fiber.Map{
			"analysis": analysis,
			"findings": findings,
		})
	}
}
