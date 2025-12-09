package server

import (
	"os"
	"os/signal"
	"syscall"

	"gorm.io/gorm"

	"mws-ai/internal/config"
	"mws-ai/internal/router"
	"mws-ai/pkg/logger"
)

func Run(cfg *config.Config, db *gorm.DB) {
	app := router.Setup(cfg, db)

	errChan := make(chan error)

	go func() {
		logger.Log.Info().Msg("Starting HTTP server on :" + cfg.ServerPort)
		if err := app.Listen(":" + cfg.ServerPort); err != nil {
			errChan <- err
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		logger.Log.Warn().Str("signal", sig.String()).Msg("Graceful shutdown triggered")
	case err := <-errChan:
		logger.Log.Error().Err(err).Msg("Server error")
	}

	if err := app.Shutdown(); err != nil {
		logger.Log.Error().Err(err).Msg("Error during server shutdown")
	}

	logger.Log.Info().Msg("Server stopped")

	sqlDB, _ := db.DB()
	_ = sqlDB.Close()
	logger.Log.Info().Msg("Database connection closed")

}
