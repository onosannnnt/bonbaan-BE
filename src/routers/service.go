package router

import (
	"github.com/gofiber/fiber/v2"
	attachmentAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/attachment" // if you have one
	serviceAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/service"
	attachmentUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/attachment"
	serviceUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/service"
	"gorm.io/gorm"
)

func InitServiceRouter(app *fiber.App, db *gorm.DB) {

    serviceRepo := serviceAdapter.NewServiceDriver(db)
    svcUsecase := serviceUsecase.NewServiceUsecase(serviceRepo)

    // Initialize the attachment repository and use case.
    attachmentRepo := attachmentAdapter.NewAttachmentDriver(db)
    attUsecase := attachmentUsecase.NewAttachmentService(attachmentRepo)

    serviceHandler := serviceAdapter.NewServiceHandler(svcUsecase, attUsecase)

    ser := app.Group("/services")
    ser.Get("/", serviceHandler.GetAllServices)
    ser.Get("/:id", serviceHandler.GetByServiceID)
    ser.Get("/:id/packages", serviceHandler.GetPackagesbyServiceID)

    // protected := ser.Group("/")
    // protected.Use(middleware.IsAuth)

    // admin := protected.Group("/")
    // admin.Use(middleware.IsAdmin)

    ser.Post("/", serviceHandler.CreateService)
    ser.Patch("/:id", serviceHandler.UpdateService)
    ser.Delete("/:id", serviceHandler.DeleteService)
}