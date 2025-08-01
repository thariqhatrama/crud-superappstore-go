package handler

import (
	"strconv"

	"FinalTask/internal/middleware"
	"FinalTask/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	ProductService service.ProductService
}

func NewProductHandler(r fiber.Router, prodService service.ProductService) {
	h := &ProductHandler{ProductService: prodService}
	group := r.Group("/products", middleware.JWTProtected())
	group.Post("", h.CreateProduct)
	group.Get("", h.ListProduct)
	group.Get("/:id", h.GetProduct)
	group.Put("/:id", h.UpdateProduct)
	group.Delete("/:id", h.DeleteProduct)
	group.Post("/:id/upload", h.UploadProductImage)
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var req service.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	prod, err := h.ProductService.Create(c.Context(), userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(prod)
}

func (h *ProductHandler) ListProduct(c *fiber.Ctx) error {
	qs := c.Queries()
	list, err := h.ProductService.List(c.Context(), qs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}

func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	// Konversi ID
	idParam := c.Params("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}
	id := uint(id64)

	prod, err := h.ProductService.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(prod)
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	// Konversi ID
	idParam := c.Params("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}
	id := uint(id64)

	var req service.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	updated, err := h.ProductService.Update(c.Context(), userID, id, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(updated)
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	// Konversi ID
	idParam := c.Params("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}
	id := uint(id64)

	if err := h.ProductService.Delete(c.Context(), userID, id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ProductHandler) UploadProductImage(c *fiber.Ctx) error {
	// Konversi ID
	idParam := c.Params("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}
	id := uint(id64)

	file, err := c.FormFile("file")
	if err != nil {
		return fiber.ErrBadRequest
	}
	url, err := h.ProductService.UploadImage(c.Context(), id, file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"url": url})
}
