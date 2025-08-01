package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// AdminOnly allows access only to users with is_admin = true
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isAdmin, ok := c.Locals("is_admin").(bool)
		if !ok || !isAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied, admin only"})
		}
		return c.Next()
	}
}
