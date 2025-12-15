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

		log.Debug().Msg("get analysis request received")

		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Warn().
				Str("id_param", idStr).
				Msg("invalid analysis id")

			return fiber.ErrBadRequest
		}

		log.Debug().
			Uint("analysis_id", uint(id)).
			Msg("analysis id parsed")

		analysis, err := h.service.GetByID(uint(id))
		if err != nil {
			log.Info().
				Uint("analysis_id", uint(id)).
				Msg("analysis not found")

			return fiber.ErrNotFound
		}

		log.Info().
			Uint("analysis_id", uint(id)).
			Msg("analysis returned")

		return c.JSON(analysis)
	}
}
