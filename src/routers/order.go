package router

import (
	"github.com/gofiber/fiber/v2"
	orderAdepter "github.com/onosannnnt/bonbaan-BE/src/adepters/order"
	serviceAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/service"
	statusAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/status"
	orderUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/order"
	statusUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/status"
	"gorm.io/gorm"
)

func InitOrderRouter(app *fiber.App, db *gorm.DB) {

	statusRepo := statusAdapter.NewStatusDriver(db)
	statusUsecase := statusUsecase.NewStatusService(statusRepo)

	orderRepo := orderAdepter.NewOrderDriver(db, statusUsecase)
	serviceRepo := serviceAdapter.NewServiceDriver(db)
	orderUsecase := orderUsecase.NewOrderService(orderRepo, serviceRepo, statusRepo)
	orderHandler := orderAdepter.NewOrderHandler(orderUsecase)

	order := app.Group("/orders")
	order.Post("/", orderHandler.Insert)
	order.Get("/", orderHandler.GetAll)
	order.Get("/:id", orderHandler.GetByID)
	order.Patch("/:id", orderHandler.Update)
	order.Delete("/:id", orderHandler.Delete)
	order.Post("/webhook", orderHandler.Hook)
	order.Post("/cancel/:id", orderHandler.CancleOrder)
	order.Post("/accept/:id", orderHandler.AcceptOrder)
	order.Post("/confirm", orderHandler.ConfirmOrder)
}
