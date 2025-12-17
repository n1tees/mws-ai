package analysis

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"mws-ai/internal/services"
	"mws-ai/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type UploadHandler struct {
	service   *services.AnalysisService
	uploadDir string
}

func NewUploadHandler(service *services.AnalysisService, uploadDir string) *UploadHandler {
	return &UploadHandler{
		service:   service,
		uploadDir: uploadDir,
	}
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
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /analysis/upload [post]
func (h *UploadHandler) Upload() fiber.Handler {
	return func(c *fiber.Ctx) error {

		log := logger.Log.With().
			Str("handler", "analysis.upload").
			Str("path", c.Path()).
			Logger()

		log.Debug().Msg("upload request received")

		userID := c.Locals("user_id")
		if userID == nil {
			log.Warn().Msg("unauthorized upload attempt")
			return fiber.ErrUnauthorized
		}
		uid := userID.(uint)

		file, err := c.FormFile("file")
		if err != nil {
			log.Warn().Err(err).Msg("file not provided")
			return fiber.NewError(fiber.StatusBadRequest, "file is required")
		}

		log.Debug().
			Uint("user_id", uid).
			Str("filename", file.Filename).
			Int64("size", file.Size).
			Msg("file received")

		// 🔥 гарантируем наличие директории
		if err := os.MkdirAll(h.uploadDir, 0755); err != nil {
			log.Error().
				Err(err).
				Str("upload_dir", h.uploadDir).
				Msg("failed to create upload directory")

			return fiber.NewError(
				fiber.StatusInternalServerError,
				"failed to prepare upload directory",
			)
		}

		filePath := filepath.Join(
			h.uploadDir,
			fmt.Sprintf("%d_%d_%s", uid, time.Now().Unix(), file.Filename),
		)

		log.Debug().
			Str("file_path", filePath).
			Msg("generated file path")

		if err := c.SaveFile(file, filePath); err != nil {
			log.Error().
				Err(err).
				Str("file_path", filePath).
				Msg("failed to save uploaded file")

			return fiber.NewError(
				fiber.StatusInternalServerError,
				"cannot save file",
			)
		}

		log.Info().
			Uint("user_id", uid).
			Str("file_path", filePath).
			Msg("file saved successfully")

		analysis, err := h.service.Upload(uid, filePath)
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to create analysis")

			return fiber.NewError(
				fiber.StatusInternalServerError,
				err.Error(),
			)
		}

		return c.JSON(fiber.Map{
			"analysis_id": analysis.ID,
			"status":      "uploaded",
		})
	}
}
