package router

import (
	"github.com/gofiber/fiber/v2"
	orderAdepter "github.com/onosannnnt/bonbaan-BE/src/adepters/order"
	orderUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/order"
	"gorm.io/gorm"
)

func InitOrderRouter(app *fiber.App, db *gorm.DB) {

	orderRepo := orderAdepter.NewOrderDriver(db)
	orderUsecase := orderUsecase.NewOrderService(orderRepo)
	orderHandler := orderAdepter.NewOrderHandler(orderUsecase)

	order := app.Group("/orders")
	order.Post("/", orderHandler.Insert)
	order.Get("/", orderHandler.GetAll)
	order.Get("/:id", orderHandler.GetByID)
	order.Put("/:id", orderHandler.Update)
	order.Delete("/:id", orderHandler.Delete)

}
