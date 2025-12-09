package auth

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: регистрация пользователя
		return c.JSON(fiber.Map{"msg": "register endpoint"})
	}
}
