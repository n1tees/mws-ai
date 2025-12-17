package middleware

import (
	"mws-ai/internal/services"
	"mws-ai/pkg/jwt"
	"mws-ai/pkg/logger"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(
	jwtManager *jwt.JWTManager,
	apiKeys *services.APIKeyService,
) fiber.Handler {
	return func(c *fiber.Ctx) error {

		log := logger.Log.With().
			Str("component", "auth").
			Str("middleware", "Auth").
			Str("path", c.Path()).
			Logger()

		authHeader := c.Get("Authorization")
		apiKey := c.Get("X-API-Key")

		if authHeader != "" && apiKey != "" {
			log.Warn().Msg("both JWT and API key provided")
			return fiber.ErrUnauthorized
		}

		//API KEY
		if apiKey != "" {
			userID, err := apiKeys.Validate(apiKey)
			if err != nil {
				log.Info().Err(err).Msg("API key validation failed")
				return fiber.ErrUnauthorized
			}

			log.Debug().
				Uint("user_id", userID).
				Msg("authorized via API key")

			c.Locals("user_id", userID)
			c.Locals("auth_type", "api_key")
			return c.Next()
		}

		//JWT
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := jwtManager.Parse(tokenStr)
			if err != nil {
				log.Info().Err(err).Msg("JWT validation failed")
				return fiber.ErrUnauthorized
			}

			userIDRaw := claims["user_id"]
			userID := uint(userIDRaw.(float64))

			log.Debug().
				Uint("user_id", userID).
				Msg("authorized via JWT")

			c.Locals("user_id", userID)
			c.Locals("auth_type", "jwt")
			return c.Next()
		}

		log.Debug().Msg("no auth provided")
		return fiber.ErrUnauthorized
	}
}
