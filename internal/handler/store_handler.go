package handler

import (
	"FinalTask/internal/middleware"
	"FinalTask/internal/service"

	"github.com/gofiber/fiber/v2"
)

type StoreHandler struct {
	StoreService service.StoreService
}

func NewStoreHandler(r fiber.Router, storeService service.StoreService) {
	h := &StoreHandler{StoreService: storeService}
	r.Get("/store", middleware.JWTProtected(), h.GetStore)
	r.Put("/store", middleware.JWTProtected(), h.UpdateStore)
}

func (h *StoreHandler) GetStore(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	store, err := h.StoreService.GetByUser(c.Context(), userID)
	if err != nil {
		return fiber.ErrNotFound
	}
	return c.JSON(store)
}

func (h *StoreHandler) UpdateStore(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var req service.UpdateStoreRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	updated, err := h.StoreService.Update(c.Context(), userID, req)
	if err != nil {
		return fiber.ErrBadRequest
	}
	return c.JSON(updated)
}
