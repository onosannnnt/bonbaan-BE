package utils

import "github.com/gofiber/fiber/v2"

func ResponseJSON(c *fiber.Ctx, statusCode int, message string, err error, data any) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"message": message,
		"error":   err.Error(),
		"data":    data,
	})
}
