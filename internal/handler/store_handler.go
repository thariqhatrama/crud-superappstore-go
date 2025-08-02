package handler

import (
	"strconv"

	"FinalTask/internal/middleware"
	"FinalTask/internal/service"

	"github.com/gofiber/fiber/v2"
)

type StoreHandler struct {
	StoreService service.StoreService
}

func NewStoreHandler(r fiber.Router, storeService service.StoreService) {
	h := &StoreHandler{StoreService: storeService}

	// --- My store endpoints (user must be logged in) ---
	storeGroup := r.Group("/store", middleware.JWTProtected())
	storeGroup.Get("", h.GetMyStore)    // GET  /store
	storeGroup.Put("", h.UpdateMyStore) // PUT  /store

	// --- Public/Admin endpoints under /stores ---
	storesGroup := r.Group("/stores", middleware.JWTProtected())
	// AdminOnly: only admin can list all or get any store
	storesGroup.Use(middleware.AdminOnly())
	storesGroup.Get("", h.GetAllStores)     // GET  /stores
	storesGroup.Get("/:id", h.GetStoreByID) // GET  /stores/:id
}

func (h *StoreHandler) GetMyStore(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	store, err := h.StoreService.GetByUser(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "store not found",
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"store": store,
		},
	})
}

func (h *StoreHandler) UpdateMyStore(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var req service.UpdateStoreRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "invalid request payload",
		})
	}
	updated, err := h.StoreService.Update(c.Context(), userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"store": updated,
		},
	})
}

func (h *StoreHandler) GetAllStores(c *fiber.Ctx) error {
	list, err := h.StoreService.ListAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to retrieve stores",
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"stores": list,
		},
	})
}

func (h *StoreHandler) GetStoreByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "invalid store ID",
		})
	}
	id := uint(id64)

	store, err := h.StoreService.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "store not found",
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"store": store,
		},
	})
}
