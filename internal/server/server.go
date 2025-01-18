package server

import (
	"github.com/gofiber/fiber/v2"

	"bonbaan-BE/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "bonbaan-BE",
			AppName:      "bonbaan-BE",
		}),

		db: database.New(),
	}

	return server
}
