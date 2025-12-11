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
