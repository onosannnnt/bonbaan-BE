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
	orderUsecase := orderUsecase.NewOrderService(orderRepo, serviceRepo, statusRepo, db)
	orderHandler := orderAdepter.NewOrderHandler(orderUsecase)

	order := app.Group("/orders")
	order.Post("/", orderHandler.Insert)
	order.Get("/", orderHandler.GetAll)
	order.Get("/:id", orderHandler.GetByID)
	order.Post("/webhook", orderHandler.Hook)

	protected := order.Group("/")
	protected.Use(middleware.IsAuth)
	protected.Patch("/:id", orderHandler.Update)
	protected.Delete("/:id", orderHandler.Delete)

	protected.Post("/:id/cancel", orderHandler.CancelOrder)
	protected.Post("/:id/submit", orderHandler.SubmitOrder)
	protected.Post("/custom-order", orderHandler.InsertCustomOrder)
	protected.Post("/:id/approve", orderHandler.ApproveOrder)
	protected.Post("/:id/complete", orderHandler.CompleteOrder)

	admin := protected.Group("/")
	admin.Use(middleware.IsAdmin)

	admin.Post("/:id/accept", orderHandler.AcceptOrder)

}
