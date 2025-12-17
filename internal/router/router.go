package router

import (
	_ "mws-ai/docs/swagger"

	swagger "github.com/gofiber/swagger"

	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"mws-ai/internal/config"
	"mws-ai/internal/router/middleware"
	"mws-ai/internal/services/clients"

	authMiddleware "mws-ai/internal/auth/middleware"
	analysisHandlers "mws-ai/internal/handlers/analysis"
	authHandlers "mws-ai/internal/handlers/auth"
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

	// INIT PARSER
	parser := sarif.NewParser()

	// INIT CLIENTS 4 PIPELINE
	heuristicClient := clients.NewHeuristicClient(cfg.HeuristicURL)
	mlClient := clients.NewMLClient(cfg.MLURL)
	llmClient := clients.NewLLMClient(cfg.LLMURL)
	// INIT PIPELINE EXECUTOR
	pipeline := services.NewPipeline(heuristicClient, mlClient, llmClient, findingRepo)

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
	uploadHandler := analysisHandlers.NewUploadHandler(analysisService, cfg.UploadDir)

	// ROUTER STRUCTURE
	api := app.Group("/api")

	// HEALTH CHECKPOINT
	api.Get("/health", healthHandlers.HealthHandler())

	// SWAGGER
	api.Get("/swagger/*", swagger.HandlerDefault)

	// AUTH ROUTES
	authGroup := api.Group("/auth")
	{
		authGroup.Post("/register", authHandler.Register())
		authGroup.Post("/login", authHandler.Login())

		// выдача API ключа
		authGroup.Post("/api-key", authMiddleware.JWTMiddleware(jwtManager), apiKeyHandler.CreateAPIKey())
	}

	// ANALYSIS ROUTES (protected)
	analysisGroup := api.Group("/analysis", middleware.AuthMiddleware(jwtManager, apiKeyService))
	{
		analysisGroup.Post("/upload", uploadHandler.Upload())
		analysisGroup.Get("/:id", analysisHandler.Get())
		analysisGroup.Get("/", analysisHandler.List())
	}

	return app
}
