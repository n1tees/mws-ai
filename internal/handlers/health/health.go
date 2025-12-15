package health

import (
	"mws-ai/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// Health godoc
// @Summary Проверка состояния сервера
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {

		logger.Log.Debug().
			Str("handler", "health").
			Str("path", c.Path()).
			Msg("health check requested")

		return c.JSON(fiber.Map{"status": "ok"})
	}
}
