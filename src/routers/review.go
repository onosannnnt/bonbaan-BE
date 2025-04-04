package router

import (
	"github.com/gofiber/fiber/v2"
	orderAdepter "github.com/onosannnnt/bonbaan-BE/src/adepters/order"
	reviewAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/review"
	statusAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/status"
	reviewUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/review"
	"github.com/onosannnnt/bonbaan-BE/src/utils/middleware"

	"gorm.io/gorm"
)

func InitReviewRouter(app *fiber.App, db *gorm.DB) {

	orderRepo := orderAdepter.NewOrderDriver(db, nil)
	statusRepo := statusAdapter.NewStatusDriver(db)

	reviewRepo := reviewAdapter.NewReviewDriver(db)
	reviewUsecase := reviewUsecase.NewReviewService(reviewRepo, orderRepo, statusRepo)
	reviewHandler := reviewAdapter.NewReviewHandler(reviewUsecase)

	review := app.Group("/reviews")
	review.Get("/", reviewHandler.GetAll)
	review.Get("/:id", reviewHandler.GetByID)

	protected := review.Group("/")
	protected.Use(middleware.IsAuth)
	protected.Post("/", reviewHandler.Insert)
	protected.Patch("/:id", reviewHandler.Update)
	protected.Delete("/:id", reviewHandler.Delete)

	// protected := review.Group("/protected")
	// protected.Use(middleware.IsAuth)

	// admin := protected.Group("/admin")
	// admin.Use(middleware.IsAdmin)

}
