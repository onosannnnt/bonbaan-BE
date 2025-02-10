package router

import (
	"github.com/gofiber/fiber/v2"
	attachmentAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/attachment"
	attachmentUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/attachment"
	"gorm.io/gorm"
)

func InitAttachmentRouter(app *fiber.App, db *gorm.DB) {
    attachmentRepo := attachmentAdapter.NewAttachmentDriver(db)
    usecase := attachmentUsecase.NewAttachmentService(attachmentRepo)
    handler := attachmentAdapter.NewAttachmentHandler(usecase)

    attachmentRouter := app.Group("/attachments")

    attachmentRouter.Get("/", handler.GetAllAttachment)
    attachmentRouter.Get("/:id", handler.GetAttachmentByID)
    attachmentRouter.Get("/service/:id", handler.GetAttachmentByServiceID)

    attachmentRouter.Post("/", handler.CreateAttachment)
    attachmentRouter.Patch("/:id", handler.UpdateAttachment)
    attachmentRouter.Delete("/:id", handler.DeleteAttachment)
}
