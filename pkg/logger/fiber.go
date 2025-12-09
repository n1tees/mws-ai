package logger

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func FiberLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start).Milliseconds()

		Log.Info().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Int64("duration_ms", duration).
			Str("ip", c.IP()).
			Msg("incoming request")

		return err
	}
}
