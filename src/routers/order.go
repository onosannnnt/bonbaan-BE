package router

import (
	"github.com/gofiber/fiber/v2"
	notificationAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/notification"
	orderAdepter "github.com/onosannnnt/bonbaan-BE/src/adepters/order"
	orderTypeAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/order_type"
	packageAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/package"
	serviceAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/service"
	statusAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/status"
	vowrecordAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/vow_record"
	orderUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/order"
	statusUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/status"
	vowRecordUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/vow_record"
	"github.com/onosannnnt/bonbaan-BE/src/utils/middleware"
	"gorm.io/gorm"
)

func InitOrderRouter(app *fiber.App, db *gorm.DB) {

	statusRepo := statusAdapter.NewStatusDriver(db)
	statusUsecase := statusUsecase.NewStatusService(statusRepo)

	vowRecordRepo := vowrecordAdapter.NewVowRecordDriver(db)
	vowRecordUsecase := vowRecordUsecase.NewVowRecordService(vowRecordRepo)
	packageRepo := packageAdapter.NewPackageDriver(db)
	orderTypeRepo := orderTypeAdapter.NewOrderTypeDriver(db)

	notificationnRepo := notificationAdapter.NewNotificationDriver(db)

	orderRepo := orderAdepter.NewOrderDriver(db, statusUsecase)
	serviceRepo := serviceAdapter.NewServiceDriver(db)
	orderUsecase := orderUsecase.NewOrderService(orderRepo, serviceRepo, statusRepo, db, packageRepo, vowRecordRepo, orderTypeRepo, notificationnRepo)
	orderHandler := orderAdepter.NewOrderHandler(orderUsecase, *vowRecordUsecase)

	order := app.Group("/orders")

	order.Get("/", orderHandler.GetAll)
	order.Get("/:id", orderHandler.GetByID)
	order.Post("/webhook", orderHandler.Hook)

	protected := order.Group("/")
	protected.Use(middleware.IsAuth)
	protected.Post("/", orderHandler.Insert)
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
