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

// CreateProduct handles POST /products
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var req service.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request payload",
		})
	}
	prod, err := h.ProductService.Create(c.Context(), userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"product": prod,
		},
	})
}

// ListProduct handles GET /products
func (h *ProductHandler) ListProduct(c *fiber.Ctx) error {
	qs := c.Queries()
	list, err := h.ProductService.List(c.Context(), qs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"products": list,
		},
	})
}

// GetProduct handles GET /products/:id
func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid product ID",
		})
	}
	id := uint(id64)

	prod, err := h.ProductService.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"product": prod,
		},
	})
}

// UpdateProduct handles PUT /products/:id
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	idParam := c.Params("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid product ID",
		})
	}
	id := uint(id64)

	var req service.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request payload",
		})
	}
	updated, err := h.ProductService.Update(c.Context(), userID, id, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"product": updated,
		},
	})
}

// DeleteProduct handles DELETE /products/:id
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	idParam := c.Params("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid product ID",
		})
	}
	id := uint(id64)

	if err := h.ProductService.Delete(c.Context(), userID, id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusNoContent).JSON(nil)
}

// UploadProductImage handles POST /products/:id/upload
func (h *ProductHandler) UploadProductImage(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid product ID",
		})
	}
	id := uint(id64)

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "File is required",
		})
	}
	url, err := h.ProductService.UploadImage(c.Context(), id, file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"url": url,
		},
	})
}
