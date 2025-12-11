package auth

import (
	"mws-ai/internal/dto"
	"mws-ai/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: service}
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Создаёт нового пользователя по email и паролю
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body dto.RegisterRequest true "Данные для регистрации"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req dto.RegisterRequest

		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		user, err := h.authService.Register(req)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return c.JSON(user)
	}
}

// Login godoc
// @Summary Авторизация пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body dto.LoginRequest true "Email и пароль"
// @Success 200 {object} dto.LoginResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req dto.LoginRequest

		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		resp, err := h.authService.Login(req)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		return c.JSON(resp)
	}
}
