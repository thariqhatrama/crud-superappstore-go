package handler

import (
	"strconv"

	"FinalTask/internal/middleware"
	"FinalTask/internal/service"

	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	TrxService service.TransactionService
}

func NewTransactionHandler(r fiber.Router, trxService service.TransactionService) {
	h := &TransactionHandler{TrxService: trxService}
	group := r.Group("/transactions", middleware.JWTProtected())
	group.Post("", h.CreateTransaction)
	group.Get("", h.ListTransactions)
	group.Get("/:id", h.GetTransaction)
}

func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var req service.CreateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	trx, err := h.TrxService.Create(c.Context(), userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(trx)
}

func (h *TransactionHandler) ListTransactions(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	qs := c.Queries()
	list, err := h.TrxService.List(c.Context(), userID, qs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}

func (h *TransactionHandler) GetTransaction(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	// Konversi ID dari string ke uint
	idParam := c.Params("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid transaction ID"})
	}
	id := uint(id64)

	trx, err := h.TrxService.GetByID(c.Context(), userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(trx)
}
