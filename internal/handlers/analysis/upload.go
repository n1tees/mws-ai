package analysis

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UploadHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: загрузка SARIF + запуск пайплайна
		return c.JSON(fiber.Map{"msg": "upload endpoint"})
	}
}
