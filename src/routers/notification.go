package router

import (
	"github.com/gofiber/fiber/v2"
	notificationAdepter "github.com/onosannnnt/bonbaan-BE/src/adepters/notification"
	notificationUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/notification"
	"gorm.io/gorm"
)

func InitNotificationRouter(app *fiber.App, db *gorm.DB) {
	Repo := notificationAdepter.NewNotificationDriver(db)
	usecase := notificationUsecase.NewNotificationService(Repo)
	handler := notificationAdepter.NewNotificationHandler(usecase)

	notificationRouter := app.Group("/notifications")
	notificationRouter.Get("/", handler.GetAll)
	notificationRouter.Get("/:id", handler.GetByID)

	notificationRouter.Post("/", handler.Insert)
	notificationRouter.Patch("/:id", handler.Update)
	notificationRouter.Delete("/:id", handler.Delete)
	notificationRouter.Patch("/:id/read", handler.Read)
}
