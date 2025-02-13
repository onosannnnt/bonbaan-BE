package router

import (
	"github.com/gofiber/fiber/v2"
	orderAdepter "github.com/onosannnnt/bonbaan-BE/src/adepters/order"
	serviceAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/service"
	orderUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/order"
	"gorm.io/gorm"
)

func InitOrderRouter(app *fiber.App, db *gorm.DB) {

	orderRepo := orderAdepter.NewOrderDriver(db)
	serviceRepo := serviceAdapter.NewServiceDriver(db)
	orderUsecase := orderUsecase.NewOrderService(orderRepo, serviceRepo)
	orderHandler := orderAdepter.NewOrderHandler(orderUsecase)

	order := app.Group("/orders")
	order.Post("/", orderHandler.Insert)
	order.Get("/", orderHandler.GetAll)
	order.Get("/:id", orderHandler.GetByID)
	order.Patch("/:id", orderHandler.Update)
	order.Delete("/:id", orderHandler.Delete)
	order.Post("/webhook", orderHandler.Hook)

}
