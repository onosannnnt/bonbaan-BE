package router

import (
	"github.com/gofiber/fiber/v2"
	categoryAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/category"
	categoryUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/category"
	"gorm.io/gorm"
)

func InitCategoryRouter(app *fiber.App, db *gorm.DB) {
	categoryRepo := categoryAdapter.NewCategoryDriver(db)
	usecase := categoryUsecase.NewCategoryUsecase(categoryRepo)
	handler := categoryAdapter.NewHandler(usecase)

	categoryRouter := app.Group("/categories")


	categoryRouter.Get("/", handler.GetAllCategory)
	categoryRouter.Get("/:id", handler.GetByCategoryID)
	
	categoryRouter.Post("/", handler.CreateCategory)
	categoryRouter.Put("/:id", handler.UpdateCategory)
	categoryRouter.Delete("/:id", handler.DeleteCategory)


	// service_category := categoryRouter.Group("/service") // Define a sub-group for service-related routes
	// service_category.Post("/", handler.AddServiceToCategory) // add service to category
	// service_category.Delete("/", handler.RemoveServiceFromCategory) // remove service from category
// 	service_category.Get("/getservices/:id", handler.GetServicesByCategoryID) // get services by category ID
// 	service_category.Get("/getcategories/:id", handler.GetCategoriesByServiceID)//  get categories by service ID
	
}
