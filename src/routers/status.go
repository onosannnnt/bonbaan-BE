package router

import (
	"github.com/gofiber/fiber/v2"
	statusAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/status"
	statusUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/status"
	"github.com/onosannnnt/bonbaan-BE/src/utils/middleware"
	"gorm.io/gorm"
)

func InitStatusRouter(app *fiber.App, db *gorm.DB) {
	statusRepo := statusAdapter.NewStatusDriver(db)
	statusUsecase := statusUsecase.NewStatusService(statusRepo)
	statusHandler := statusAdapter.NewStatusHandler(statusUsecase)

	status := app.Group("/statuses")
	status.Get("/", statusHandler.GetAllStatus)
	status.Get("/:id", statusHandler.GetStatusByID)
	status.Get("/name/:name", statusHandler.GetStatusByName)

	protected := status.Group("/")
	protected.Use(middleware.IsAuth)

	admin := protected.Group("/")
	admin.Use(middleware.IsAdmin)
	admin.Post("/", statusHandler.InsertStatus)
	admin.Patch("/:id", statusHandler.UpdateStatus)
	admin.Delete("/:id", statusHandler.DeleteStatus)

}
