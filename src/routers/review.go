package router

import (
	"github.com/gofiber/fiber/v2"
	reviewAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/review"
	reviewUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/review"

	"gorm.io/gorm"
)

func ReviewRouter(app *fiber.App, db *gorm.DB) {

	reviewRepo := reviewAdapter.NewReviewDriver(db)
	reviewUsecase := reviewUsecase.NewReviewService(reviewRepo)
	reviewHandler := reviewAdapter.NewReviewHandler(reviewUsecase)

	review := app.Group("/reviews")
	review.Post("/", reviewHandler.Insert)


	// review.Delete("/:id",reviewHandler.Delete)

	// protected := review.Group("/protected")
	// protected.Use(middleware.IsAuth)

	// admin := protected.Group("/admin")
	// admin.Use(middleware.IsAdmin)

	
	

}