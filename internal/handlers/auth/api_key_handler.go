package auth

import (
	"mws-ai/internal/services"
	"mws-ai/pkg/logger"
	"time"

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

		log := logger.Log.With().
			Str("component", "auth").
			Str("handler", "CreateAPIKey").
			Logger()

		userIDRaw := c.Locals("user_id")
		userID, ok := userIDRaw.(uint)
		if !ok {
			log.Warn().
				Interface("user_id_raw", userIDRaw).
				Msg("user_id missing or invalid in context")

			return fiber.ErrUnauthorized
		}

		//ADMIN ONLY
		if userID != 1 {
			log.Warn().
				Uint("user_id", userID).
				Msg("non-admin attempted to create API key")

			return fiber.ErrForbidden
		}

		log.Info().
			Uint("user_id", userID).
			Msg("admin requested API key generation")

		// TTL
		rawKey, err := h.apiKeys.Generate(
			userID,
			"service",
			7*24*time.Hour,
		)
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to generate API key")

			return fiber.NewError(
				fiber.StatusInternalServerError,
				"failed to generate api key",
			)
		}

		log.Info().
			Msg("API key generated successfully")

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"api_key": rawKey,
			"type":    "service",
			"expires": "7d",
			"note":    "Store this key securely. It will not be shown again.",
		})
	}
}
