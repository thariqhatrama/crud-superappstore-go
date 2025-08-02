package handler

import (
	"FinalTask/internal/middleware"
	"FinalTask/internal/service"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserService service.UserService
}

func NewUserHandler(r fiber.Router, userService service.UserService) {
	h := &UserHandler{UserService: userService}
	r.Get("/me", middleware.JWTProtected(), h.GetProfile)
	r.Put("/me", middleware.JWTProtected(), h.UpdateProfile)
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	user, err := h.UserService.GetByID(c.Context(), userID)
	if err != nil {
		return fiber.ErrNotFound
	}
	return c.JSON(user)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req service.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "invalid request payload",
		})
	}

	updated, err := h.UserService.Update(c.Context(), userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   fiber.Map{"user": updated},
	})
}
