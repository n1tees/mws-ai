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

		status := c.Response().StatusCode()

		if err != nil {
			if fe, ok := err.(*fiber.Error); ok {
				status = fe.Code
			}
		}

		duration := time.Since(start).Milliseconds()

		logger.Log.Info().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", status).
			Int64("duration_ms", duration).
			Str("ip", c.IP()).
			Msg("incoming request")

		return err
	}
}
