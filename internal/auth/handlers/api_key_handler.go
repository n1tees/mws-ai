package auth

import (
	"mws-ai/internal/services"
	"mws-ai/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type APIKeyHandler struct {
	apiKeys *services.APIKeyService
}

func NewAPIKeyHandler(apiKeys *services.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{apiKeys: apiKeys}
}

// CreateAPIKey godoc
// @Summary Создание API-ключа
// @Description Генерирует новый API-ключ для текущего пользователя
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 201 {object} dto.CreateAPIKeyResponse
// @Failure 401 {object} dto.ErrorResponse "Неавторизован"
// @Router /auth/api-key [post]
func (h *APIKeyHandler) CreateAPIKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userIDRaw := c.Locals("user_id")
		userID, ok := userIDRaw.(uint)
		if !ok {
			logger.Log.Warn().
				Str("component", "auth").
				Str("handler", "CreateAPIKey").
				Interface("user_id_raw", userIDRaw).
				Msg("user_id missing or invalid in context")

			return fiber.ErrUnauthorized
		}

		logger.Log.Debug().
			Str("component", "auth").
			Str("handler", "CreateAPIKey").
			Uint("user_id", userID).
			Msg("API key generation requested")

		rawKey, err := h.apiKeys.Generate(userID)
		if err != nil {
			logger.Log.Error().
				Err(err).
				Str("component", "auth").
				Str("handler", "CreateAPIKey").
				Uint("user_id", userID).
				Msg("failed to generate API key")

			return fiber.ErrInternalServerError
		}

		logger.Log.Info().
			Str("component", "auth").
			Str("handler", "CreateAPIKey").
			Uint("user_id", userID).
			Msg("API key generated successfully")

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"api_key": rawKey,
		})
	}
}
