package repository

import (
	"errors"

	"mws-ai/internal/models"
	"mws-ai/pkg/logger"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
	Update(user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		logger.Log.Error().
			Str("repo", "user").
			Str("method", "Create").
			Str("email", user.Email).
			Err(err).
			Msg("failed to create user")

		return err
	}

	logger.Log.Debug().
		Str("repo", "user").
		Str("method", "Create").
		Uint("user_id", user.ID).
		Str("email", user.Email).
		Msg("user created")

	return nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User

	err := r.db.
		Where("email = ?", email).
		First(&user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		logger.Log.Error().
			Str("repo", "user").
			Str("method", "FindByEmail").
			Str("email", email).
			Err(err).
			Msg("failed to find user by email")

		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User

	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		logger.Log.Error().
			Str("repo", "user").
			Str("method", "FindByID").
			Uint("user_id", id).
			Err(err).
			Msg("failed to find user by id")

		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	res := r.db.Save(user)

	if res.Error != nil {
		logger.Log.Error().
			Str("repo", "user").
			Str("method", "Update").
			Uint("user_id", user.ID).
			Err(res.Error).
			Msg("failed to update user")

		return res.Error
	}

	if res.RowsAffected == 0 {
		logger.Log.Debug().
			Str("repo", "user").
			Str("method", "Update").
			Uint("user_id", user.ID).
			Msg("no user updated")
	}

	return nil
}
