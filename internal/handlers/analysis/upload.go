package analysis

import (
	"fmt"
	"time"

	"mws-ai/internal/services"

	"github.com/gofiber/fiber/v2"
)

type UploadHandler struct {
	service *services.AnalysisService
}

func NewUploadHandler(service *services.AnalysisService) *UploadHandler {
	return &UploadHandler{service: service}
}

func (h *UploadHandler) Upload() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Авторизация: получаем user_id
		userID := c.Locals("user_id")
		if userID == nil {
			return fiber.ErrUnauthorized
		}
		uid := userID.(uint)

		// Получаем файл из multipart/form-data
		file, err := c.FormFile("file")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "file is required")
		}

		// Генерируем путь
		filePath := fmt.Sprintf("uploads/%d_%d_%s",
			uid,
			time.Now().Unix(),
			file.Filename,
		)

		// сохраняем файл
		if err := c.SaveFile(file, filePath); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "cannot save file")
		}

		// создаём запись анализа (service пишет в БД)
		analysis, err := h.service.Upload(uid, filePath)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		// Возвращаем ID анализа
		return c.JSON(fiber.Map{
			"analysis_id": analysis.ID,
			"status":      "uploaded",
		})
	}
}
