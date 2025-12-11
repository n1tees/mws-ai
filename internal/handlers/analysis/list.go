package analysis

import (
	"github.com/gofiber/fiber/v2"
)

// List godoc
// @Summary Получить список анализов пользователя
// @Description Возвращает список анализов, созданных текущим пользователем
// @Tags Analysis
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.AnalysisListItem
// @Failure 401 {object} dto.ErrorResponse "Неавторизован"
// @Router /analysis [get]
func (h *AnalysisHandler) List() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uint)

		analyses, err := h.service.ListByUser(userID)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(analyses)
	}
}
