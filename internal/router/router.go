package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"mws-ai/internal/config"
	"mws-ai/internal/router/middleware"

	authHandlers "mws-ai/internal/auth/handlers"
	authMiddleware "mws-ai/internal/auth/middleware"
	analysisHandlers "mws-ai/internal/handlers/analysis"
	healthHandlers "mws-ai/internal/handlers/health"

	"mws-ai/internal/repository"
	sarif "mws-ai/internal/sarif"
	"mws-ai/internal/services"
	jwtpkg "mws-ai/pkg/jwt"
)

func Setup(cfg *config.Config, db *gorm.DB) *fiber.App {
	app := fiber.New()

	middleware.DefaultMiddleware(app)

	jwtManager := jwtpkg.NewJWTManager(
		cfg.JWTSecret,
		15*time.Minute,
		7*24*time.Hour,
	)

	// INIT REPOSITORIES
	userRepo := repository.NewUserRepository(db)
	apiKeyRepo := repository.NewAPIKeyRepository(db)
	analysisRepo := repository.NewAnalysisRepository(db)
	findingRepo := repository.NewFindingRepository(db)

	// INIT PARSER + PIPELINERSER
	parser := sarif.NewParser()
	pipeline := &services.DummyPipeline{}

	// INIT SERVICES
	authService := services.NewAuthService(userRepo, jwtManager)
	apiKeyService := services.NewAPIKeyService(apiKeyRepo)
	analysisService := services.NewAnalysisService(
		analysisRepo,
		findingRepo,
		parser,
		pipeline,
	)
	// INIT HANDLERS
	authHandler := authHandlers.NewAuthHandler(authService)
	apiKeyHandler := authHandlers.NewAPIKeyHandler(apiKeyService)

	analysisHandler := analysisHandlers.NewAnalysisHandler(analysisService)
	uploadHandler := analysisHandlers.NewUploadHandler(analysisService)

	// ROUTER STRUCTURE
	api := app.Group("/api")

	// health
	api.Get("/health", healthHandlers.HealthHandler())

	// AUTH ROUTES
	authGroup := api.Group("/auth")
	{
		authGroup.Post("/register", authHandler.Register())
		authGroup.Post("/login", authHandler.Login())

		// выдача API ключа
		authGroup.Post("/api-key", authMiddleware.JWTMiddleware(jwtManager), apiKeyHandler.CreateAPIKey())
	}

	// ANALYSIS ROUTES (protected)
	analysisGroup := api.Group("/analysis", authMiddleware.JWTMiddleware(jwtManager))
	{
		analysisGroup.Post("/upload", uploadHandler.Upload())
		analysisGroup.Get("/:id", analysisHandler.Get())
		analysisGroup.Get("/", analysisHandler.List())
	}

	return app
}
