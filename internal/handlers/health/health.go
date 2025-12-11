package health

import "github.com/gofiber/fiber/v2"

// Health godoc
// @Summary Проверка состояния сервера
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	}
}
