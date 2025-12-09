package auth

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func LoginHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: логин пользователя
		return c.JSON(fiber.Map{"msg": "login endpoint"})
	}
}
