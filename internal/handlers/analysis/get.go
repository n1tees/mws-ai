package analysis

import (
	"strconv"

	"mws-ai/internal/repository"

	"github.com/gofiber/fiber/v2"
)

type AnalysisHandler struct {
	repo repository.AnalysisRepository
}

func NewAnalysisHandler(repo repository.AnalysisRepository) *AnalysisHandler {
	return &AnalysisHandler{repo: repo}
}

func (h *AnalysisHandler) Get() fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fiber.ErrBadRequest
		}

		analysis, err := h.repo.GetByID(uint(id))
		if err != nil {
			return fiber.ErrNotFound
		}

		return c.JSON(analysis)
	}
}
