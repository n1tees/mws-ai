package analysis

import (
	"github.com/gofiber/fiber/v2"
)

func (h *AnalysisHandler) List() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uint)

		analyses, err := h.repo.ListByUser(userID)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(analyses)
	}
}
