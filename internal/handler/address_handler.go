package handler

import (
	"strconv"

	"FinalTask/internal/middleware"
	"FinalTask/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AddressHandler struct {
	AddressService service.AddressService
}

func NewAddressHandler(r fiber.Router, addressService service.AddressService) {
	h := &AddressHandler{AddressService: addressService}

	// Protected: CRUD alamat
	addrGroup := r.Group("/addresses", middleware.JWTProtected())
	addrGroup.Post("", h.CreateAddress)
	addrGroup.Get("", h.ListAddress)
	addrGroup.Get("/:id", h.GetAddressByID)
	addrGroup.Put("/:id", h.UpdateAddress)
	addrGroup.Delete("/:id", h.DeleteAddress)

	// Public: wilayah tanpa JWT
	pub := r.Group("") // temp group sama prefix r
	pub.Get("/addresses/provinces", h.GetProvinces)
	pub.Get("/addresses/provinces/:province_id", h.GetProvinceByID)
	pub.Get("/addresses/regencies/:province_id", h.GetRegenciesByProvince)
	pub.Get("/addresses/regencies/detail/:regency_id", h.GetRegencyByID)
}

// CreateAddress handles POST /addresses
func (h *AddressHandler) CreateAddress(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req service.CreateAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request payload",
		})
	}

	addr, err := h.AddressService.Create(c.Context(), userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"address": addr,
		},
	})
}

// ListAddress handles GET /addresses
func (h *AddressHandler) ListAddress(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	list, err := h.AddressService.List(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to list addresses",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"addresses": list,
		},
	})
}

// GetAddressByID handles GET /addresses/:id
func (h *AddressHandler) GetAddressByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	idParam := c.Params("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid address ID",
		})
	}
	id := uint(id64)

	addr, err := h.AddressService.GetByID(c.Context(), userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   fiber.Map{"address": addr},
	})
}

// UpdateAddress handles PUT /addresses/:id
func (h *AddressHandler) UpdateAddress(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	idParam := c.Params("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid address ID",
		})
	}
	id := uint(id64)

	var req service.CreateAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request payload",
		})
	}

	updated, err := h.AddressService.Update(c.Context(), userID, id, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"address": updated,
		},
	})
}

// DeleteAddress handles DELETE /addresses/:id
func (h *AddressHandler) DeleteAddress(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	idParam := c.Params("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid address ID",
		})
	}
	id := uint(id64)

	if err := h.AddressService.Delete(c.Context(), userID, id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"status":  "success",
		"message": "Address deleted successfully",
	})
}

// GetProvinces returns all provinces (public)
func (h *AddressHandler) GetProvinces(c *fiber.Ctx) error {
	data, err := h.AddressService.GetProvinces(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(fiber.Map{
		"status":    "success",
		"provinces": data,
	})
}

// GetProvinceByID returns one province by ID (public)
func (h *AddressHandler) GetProvinceByID(c *fiber.Ctx) error {
	provinceID := c.Params("province_id")
	raw, err := h.AddressService.GetProvinces(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	list, ok := raw.([]interface{})
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "invalid data format")
	}
	for _, item := range list {
		if m, ok := item.(map[string]interface{}); ok {
			// m["id"] is string
			if idStr, ok := m["id"].(string); ok && idStr == provinceID {
				return c.JSON(fiber.Map{
					"status":   "success",
					"province": m,
				})
			}
		}
	}
	return fiber.NewError(fiber.StatusNotFound, "province not found")
}

// GetRegenciesByProvince returns all regencies for a given province (public)
func (h *AddressHandler) GetRegenciesByProvince(c *fiber.Ctx) error {
	provinceID := c.Params("province_id")
	data, err := h.AddressService.GetRegenciesByProvince(c.Context(), provinceID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(fiber.Map{
		"status":    "success",
		"regencies": data,
	})
}

// GetRegencyByID returns one regency by its ID (public)
func (h *AddressHandler) GetRegencyByID(c *fiber.Ctx) error {
	regencyID := c.Params("regency_id")

	// search through all provinces to find the matching regency
	rawProv, err := h.AddressService.GetProvinces(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	provs, _ := rawProv.([]interface{})
	for _, item := range provs {
		if pm, ok := item.(map[string]interface{}); ok {
			if provID, ok := pm["id"].(string); ok {
				rawReg, err := h.AddressService.GetRegenciesByProvince(c.Context(), provID)
				if err != nil {
					continue
				}
				if regs, ok := rawReg.([]interface{}); ok {
					for _, ri := range regs {
						if rm, ok := ri.(map[string]interface{}); ok {
							if idStr, ok := rm["id"].(string); ok && idStr == regencyID {
								return c.JSON(fiber.Map{
									"status":  "success",
									"regency": rm,
								})
							}
						}
					}
				}
			}
		}
	}
	return fiber.NewError(fiber.StatusNotFound, "regency not found")
}
