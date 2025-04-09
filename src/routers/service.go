package router

import (
	"github.com/gofiber/fiber/v2"
	attachmentAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/attachment"
	recommendationAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/recommendation"
	serviceAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/service"
	attachmentUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/attachment"
	recommendationUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/recommendation"
	serviceUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/service"
	"github.com/onosannnnt/bonbaan-BE/src/utils/middleware"

	"gorm.io/gorm"
)

func InitServiceRouter(app *fiber.App, db *gorm.DB) {

    serviceRepo := serviceAdapter.NewServiceDriver(db)
    svcUsecase := serviceUsecase.NewServiceUsecase(serviceRepo)

    // Initialize the attachment repository and use case. 
    attachmentRepo := attachmentAdapter.NewAttachmentDriver(db)
    attUsecase := attachmentUsecase.NewAttachmentService(attachmentRepo)

    // Initialize the recommendation use case.
    recRepo := recommendationAdapter.NewRecommendationDriver(db)
    recUsecase := recommendationUsecase.NewRecommendationService(recRepo)

    // Pass the recommendation use case and DB to the handler.
    serviceHandler := serviceAdapter.NewServiceHandler(svcUsecase, attUsecase, recUsecase, db)

    ser := app.Group("/services")
    ser.Get("/", serviceHandler.GetAllServices)

    protected := ser.Group("/")
    protected.Use(middleware.IsAuth)

    // admin := protected.Group("/")
    // admin.Use(middleware.IsAdmin)

    ser.Post("/", serviceHandler.CreateService)
    ser.Patch("/:id", serviceHandler.UpdateService)
    ser.Delete("/:id", serviceHandler.DeleteService)

    // New endpoints for recommendations and bestsellers.
    protected.Get("/recommend", serviceHandler.RecommendService)
    protected.Get("/bestseller", serviceHandler.Bestseller)
    ser.Get("/:id", serviceHandler.GetByServiceID)
    ser.Get("/:id/packages", serviceHandler.GetPackagesbyServiceID)

}