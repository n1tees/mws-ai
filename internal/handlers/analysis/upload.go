package analysis

import (
	"fmt"
	"time"

	"mws-ai/internal/services"
	"mws-ai/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type UploadHandler struct {
	service *services.AnalysisService
}

func NewUploadHandler(service *services.AnalysisService) *UploadHandler {
	return &UploadHandler{service: service}
}

// Upload godoc
// @Summary Загрузить SARIF файл на анализ
// @Description Принимает SARIF JSON, создаёт Analysis и запускает pipeline обработки
// @Tags Analysis
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "SARIF файл"
// @Success 200 {object} dto.UploadAnalysisResponse
// @Failure 400 {object} dto.ErrorResponse "Некорректный файл"
// @Failure 401 {object} dto.ErrorResponse "Неавторизован"
// @Router /analysis/upload [post]
func (h *UploadHandler) Upload() fiber.Handler {
	return func(c *fiber.Ctx) error {

		log := logger.Log.With().
			Str("handler", "analysis.upload").
			Str("path", c.Path()).
			Logger()

		log.Debug().Msg("upload request received")

		// Авторизация: получаем user_id
		userID := c.Locals("user_id")
		if userID == nil {
			log.Warn().
				Msg("unauthorized upload attempt")

			return fiber.ErrUnauthorized
		}
		uid := userID.(uint)

		log.Debug().
			Uint("user_id", uid).
			Msg("user authorized")

		// Получаем файл из multipart/form-data
		file, err := c.FormFile("file")
		if err != nil {
			log.Warn().
				Err(err).
				Uint("user_id", uid).
				Msg("file not provided in request")

			return fiber.NewError(fiber.StatusBadRequest, "file is required")
		}

		log.Debug().
			Uint("user_id", uid).
			Str("filename", file.Filename).
			Int64("size", file.Size).
			Msg("file received")

		// Генерируем путь
		filePath := fmt.Sprintf("uploads/%d_%d_%s",
			uid,
			time.Now().Unix(),
			file.Filename,
		)

		log.Debug().
			Str("file_path", filePath).
			Msg("generated file path")

		// сохраняем файл
		if err := c.SaveFile(file, filePath); err != nil {
			log.Error().
				Err(err).
				Uint("user_id", uid).
				Str("file_path", filePath).
				Msg("failed to save uploaded file")

			return fiber.NewError(fiber.StatusInternalServerError, "cannot save file")
		}

		log.Info().
			Uint("user_id", uid).
			Str("file_path", filePath).
			Msg("file saved successfully")

		// создаём запись анализа (service пишет в БД)
		analysis, err := h.service.Upload(uid, filePath)
		if err != nil {
			log.Error().
				Err(err).
				Uint("user_id", uid).
				Str("file_path", filePath).
				Msg("failed to create analysis")

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		log.Info().
			Uint("user_id", uid).
			Uint("analysis_id", analysis.ID).
			Msg("analysis upload accepted")

		// Возвращаем ID анализа
		return c.JSON(fiber.Map{
			"analysis_id": analysis.ID,
			"status":      "uploaded",
		})
	}
}
