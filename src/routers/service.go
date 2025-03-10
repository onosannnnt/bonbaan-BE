package router

import (
	"github.com/gofiber/fiber/v2"
	serviceAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/service"
	serviceUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/service"
	"gorm.io/gorm"
)

func InitServiceRouter(app *fiber.App, db *gorm.DB) {

	serviceRepo := serviceAdapter.NewServiceDriver(db)
	serviceUsecase := serviceUsecase.NewServiceUsecase(serviceRepo)
	serviceHandler := serviceAdapter.NewServiceHandler(serviceUsecase)

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
