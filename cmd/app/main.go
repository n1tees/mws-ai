package main

import (
	"log"

	"github.com/joho/godotenv"

	"mws-ai/internal/config"
	"mws-ai/internal/db"
	"mws-ai/pkg/logger"
)

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
	db.Migrate(database)

	// // Fiber
	// app := fiber.New()

	// // middleware для логирования HTTP запросов
	// app.Use(logger.FiberLogger())

	// // тестовый эндпоинт
	// app.Get("/ping", func(c *fiber.Ctx) error {
	// 	return c.JSON(fiber.Map{"msg": "pong"})
	// })

	// // 7. Запускаем сервер
	// logger.Log.Info().Str("port", cfg.ServerPort).Msg("Server starting")

	// if err := app.Listen(":" + cfg.ServerPort); err != nil {
	// 	logger.Log.Fatal().Err(err).Msg("Server crashed")
	//}
}
