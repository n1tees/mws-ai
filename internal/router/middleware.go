package router

import (
	"mws-ai/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func DefaultMiddleware(app *fiber.App) {
	app.Use(logger.FiberLogger())
	app.Use(recover.New())
	app.Use(cors.New())
}
