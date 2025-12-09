package analysis

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: вернуть анализ по id
		return c.JSON(fiber.Map{"msg": "get analysis"})
	}
}
