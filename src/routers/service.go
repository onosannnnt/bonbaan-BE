package router

import (
	"github.com/gofiber/fiber/v2"
	serviceAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/service"
	serviceUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/service"
	"gorm.io/gorm"
)

func ServiceRouter(app *fiber.App, db *gorm.DB) {

	serviceRepo := serviceAdapter.NewServiceDriver(db)
	serviceUsecase := serviceUsecase.NewServiceUsecase(serviceRepo)
	serviceHandler := serviceAdapter.NewServiceHandler(serviceUsecase)

	ser := app.Group("/services")
	ser.Post("/create", serviceHandler.CreateService)
	ser.Get("/", serviceHandler.GetAll)
	ser.Get("/:id", serviceHandler.GetByID)
	ser.Patch("/update/:id", serviceHandler.UpdateService) // Updated route to include ID parameter
	
}
