package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("Hello, world!")
	app := fiber.New()
	app.Get("/", Hello)

	app.Listen(":8080")
}

func Hello(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Hello, World!"})
}
