package router

import (
	"github.com/gofiber/fiber/v2"
	serviceAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/service"
	serviceUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/service"
	"github.com/onosannnnt/bonbaan-BE/src/utils/middleware"
	"gorm.io/gorm"
)

func ServiceRouter(app *fiber.App, db *gorm.DB) {

	serviceRepo := serviceAdapter.NewServiceDriver(db)
	serviceUsecase := serviceUsecase.NewServiceUsecase(serviceRepo)
	serviceHandler := serviceAdapter.NewServiceHandler(serviceUsecase)

	ser := app.Group("/services")
	ser.Get("/", serviceHandler.GetAllServices)
	ser.Get("/:id", serviceHandler.GetByServiceID)

	protected := ser.Group("/")
	protected.Use(middleware.IsAuth)

	admin := protected.Group("/")
	admin.Use(middleware.IsAdmin)
	admin.Post("/", serviceHandler.CreateService)
	admin.Patch("/:id", serviceHandler.UpdateService)

}
