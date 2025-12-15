package analysis

import (
	"mws-ai/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// List godoc
// @Summary Получить список анализов пользователя
// @Description Возвращает список анализов, созданных текущим пользователем
// @Tags Analysis
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.AnalysisListItem
// @Failure 401 {object} dto.ErrorResponse "Неавторизован"
// @Router /analysis [get]
func (h *AnalysisHandler) List() fiber.Handler {
	return func(c *fiber.Ctx) error {

		log := logger.Log.With().
			Str("handler", "analysis.list").
			Str("path", c.Path()).
			Logger()

		log.Debug().Msg("list analyses request received")

		userID := c.Locals("user_id").(uint)

		log.Debug().
			Uint("user_id", userID).
			Msg("user authorized")

		analyses, err := h.service.ListByUser(userID)
		if err != nil {
			log.Error().
				Err(err).
				Uint("user_id", userID).
				Msg("failed to list analyses")

			return fiber.ErrInternalServerError
		}

		log.Info().
			Uint("user_id", userID).
			Int("count", len(analyses)).
			Msg("analyses list returned")

		return c.JSON(analyses)
	}
}
