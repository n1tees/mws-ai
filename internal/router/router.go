package router

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"mws-ai/internal/config"
	"mws-ai/internal/handlers/analysis"
	"mws-ai/internal/handlers/auth"
	"mws-ai/internal/handlers/health"
	"mws-ai/internal/router/middleware"
)

func Setup(cfg *config.Config, db *gorm.DB) *fiber.App {
	app := fiber.New()

	middleware.DefaultMiddleware(app)

	api := app.Group("/api")

	api.Get("/health", health.HealthHandler())

	authGroup := api.Group("/auth")
	{
		authGroup.Post("/register", auth.RegisterHandler(db))
		authGroup.Post("/login", auth.LoginHandler(db))
	}

	analysisGroup := api.Group("/analysis")
	{
		analysisGroup.Post("/upload", analysis.UploadHandler(db))
		analysisGroup.Get("/:id", analysis.GetHandler(db))
		analysisGroup.Get("/", analysis.ListHandler(db))
	}

	return app
}
