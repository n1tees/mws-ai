package auth

import (
	"mws-ai/internal/services"

	"github.com/gofiber/fiber/v2"
)

type APIKeyHandler struct {
	apiKeys *services.APIKeyService
}

func NewAPIKeyHandler(apiKeys *services.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{apiKeys: apiKeys}
}

func (h *APIKeyHandler) CreateAPIKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uint)

		rawKey, err := h.apiKeys.Generate(userID)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{
			"api_key": rawKey,
		})
	}
}
