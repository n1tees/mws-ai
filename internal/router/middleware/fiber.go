package middleware

import (
	"time"

	"mws-ai/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func FiberLoggerMW() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start).Milliseconds()

		logger.Log.Info().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Int64("duration_ms", duration).
			Str("ip", c.IP()).
			Msg("incoming request")

		return err
	}
}
