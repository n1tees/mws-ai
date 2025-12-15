package auth

import (
	"strings"

	"mws-ai/pkg/jwt"
	"mws-ai/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware(jwtManager *jwt.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")

		if authHeader == "" {
			logger.Log.Debug().
				Str("component", "auth").
				Str("middleware", "JWT").
				Str("path", c.Path()).
				Msg("missing Authorization header")

			return fiber.ErrUnauthorized
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Log.Debug().
				Str("component", "auth").
				Str("middleware", "JWT").
				Str("path", c.Path()).
				Msg("invalid Authorization header format")

			return fiber.ErrUnauthorized
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwtManager.Parse(tokenStr)
		if err != nil {
			logger.Log.Info().
				Str("component", "auth").
				Str("middleware", "JWT").
				Str("path", c.Path()).
				Err(err).
				Msg("JWT validation failed")

			return fiber.ErrUnauthorized
		}

		userIDRaw, ok := claims["user_id"]
		if !ok {
			logger.Log.Warn().
				Str("component", "auth").
				Str("middleware", "JWT").
				Str("path", c.Path()).
				Msg("user_id claim missing in JWT")

			return fiber.ErrUnauthorized
		}

		userIDFloat, ok := userIDRaw.(float64)
		if !ok {
			logger.Log.Warn().
				Str("component", "auth").
				Str("middleware", "JWT").
				Interface("user_id_claim", userIDRaw).
				Msg("invalid user_id claim type")

			return fiber.ErrUnauthorized
		}

		userID := uint(userIDFloat)

		logger.Log.Debug().
			Str("component", "auth").
			Str("middleware", "JWT").
			Uint("user_id", userID).
			Str("path", c.Path()).
			Msg("JWT validated successfully")

		c.Locals("user_id", userID)

		return c.Next()
	}
}
