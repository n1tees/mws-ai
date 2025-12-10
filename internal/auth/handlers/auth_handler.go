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
