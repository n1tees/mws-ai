package analysis

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ListHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: список анализов
		return c.JSON(fiber.Map{"msg": "list analyses"})
	}
}
