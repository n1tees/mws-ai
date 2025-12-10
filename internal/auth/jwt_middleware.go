package auth

import (
	"strings"

	"mws-ai/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware(jwtManager *jwt.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return fiber.ErrUnauthorized
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwtManager.Parse(tokenStr)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		userIDFloat := claims["user_id"].(float64)
		userID := uint(userIDFloat)

		c.Locals("user_id", userID)

		return c.Next()
	}
}
