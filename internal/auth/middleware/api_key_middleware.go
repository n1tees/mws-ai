package auth

import (
	"strings"

	"mws-ai/internal/services"
	"mws-ai/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func APIKeyMiddleware(apiKeys *services.APIKeyService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		log := logger.Log.With().
			Str("component", "auth").
			Str("middleware", "APIKey").
			Str("path", c.Path()).
			Logger()

		rawKey := c.Get("X-API-Key")
		source := "x-api-key"

		if rawKey == "" {
			auth := c.Get("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				rawKey = strings.TrimPrefix(auth, "Bearer ")
				source = "authorization"
			}
		}

		if rawKey == "" {
			log.Debug().
				Msg("API key missing")

			return fiber.ErrUnauthorized
		}

		log.Debug().
			Str("source", source).
			Msg("API key provided")

		userID, err := apiKeys.Validate(rawKey)
		if err != nil {
			log.Info().
				Err(err).
				Str("source", source).
				Msg("API key validation failed")

			return fiber.ErrUnauthorized
		}

		log.Debug().
			Uint("user_id", userID).
			Str("source", source).
			Msg("API key validated successfully")

		c.Locals("user_id", userID)
		return c.Next()
	}
}
