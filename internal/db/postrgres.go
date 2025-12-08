package db

import (
	"mws-ai/internal/config"
	"mws-ai/internal/models"
	"mws-ai/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Init(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DBConnString()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}

	logger.Log.Info().Msg("Connected to PostgreSQL")

	return db, nil
}

func Migrate(db *gorm.DB) error {
	logger.Log.Info().Msg("Running DB migrations...")

	return db.AutoMigrate(
		&models.User{},
		&models.Analysis{},
		&models.Finding{},
	)
}

func Close(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to get DB SQL connection for closing")
		return
	}

	if err := sqlDB.Close(); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to close DB connection")
		return
	}

	logger.Log.Info().Msg("Database connection closed")
}
