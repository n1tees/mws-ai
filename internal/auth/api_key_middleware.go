package auth

import (
	"strings"

	"mws-ai/internal/services"

	"github.com/gofiber/fiber/v2"
)

func APIKeyMiddleware(apiKeys *services.APIKeyService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		rawKey := c.Get("X-API-Key")
		if rawKey == "" {
			auth := c.Get("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				rawKey = strings.TrimPrefix(auth, "Bearer ")
			}
		}

		if rawKey == "" {
			return fiber.ErrUnauthorized
		}

		userID, err := apiKeys.Validate(rawKey)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		c.Locals("user_id", userID)
		return c.Next()
	}
}
