package router

import (
	"github.com/gofiber/fiber/v2"
	orderAdepter "github.com/onosannnnt/bonbaan-BE/src/adepters/order"
	serviceAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/service"
	statusAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/status"
	orderUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/order"
	statusUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/status"
	"github.com/onosannnnt/bonbaan-BE/src/utils/middleware"
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

	protected := order.Group("/")
	protected.Use(middleware.IsAuth)
	protected.Patch("/:id", orderHandler.Update)
	protected.Delete("/:id", orderHandler.Delete)
	protected.Post("/webhook", orderHandler.Hook)
	protected.Post("/cancel/:id", orderHandler.CancleOrder)
	protected.Post("/submit", orderHandler.SubmitOrder)

	admin := protected.Group("/")
	admin.Use(middleware.IsAdmin)

	admin.Post("/accept/:id", orderHandler.AcceptOrder)
	admin.Post("/complete/:id", orderHandler.CompleteOrder)
}
