package main

import (
	"log"

	"github.com/joho/godotenv"

	"mws-ai/internal/config"
	"mws-ai/internal/db"
	"mws-ai/internal/server"
	"mws-ai/pkg/logger"
)

// @title MWS AI API
// @version 1.0
// @description API для анализа SARIF и определения FP/TP.
// @BasePath /api
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// config
	_ = godotenv.Load()

	cfg, warnings, err := config.Load()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	// logger
	logger.Init()
	logger.Log.Info().Msg("Logger initialized")

	for _, w := range warnings {
		logger.Log.Warn().Msg(w)
	}

	// db
	database, err := db.Init(cfg)
	if err != nil {
		log.Fatalf("BD error: %v", err)
	}

	db.Migrate(database)

	// fiber server
	server.Run(cfg, database)

}
