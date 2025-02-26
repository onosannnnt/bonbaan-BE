package router

import (
	"github.com/gofiber/fiber/v2"
	packageTypeAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/packageType"
	packageTypeUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/packageType"
	"github.com/onosannnnt/bonbaan-BE/src/utils/middleware"
	"gorm.io/gorm"
)

func InitPackageTypeRouter(app *fiber.App, db *gorm.DB) {
	
	packageTypeRepo := packageTypeAdapter.NewPackageTypeDriver(db)
	packageTypeUsecase := packageTypeUsecase.NewPackageTypeService(packageTypeRepo)
	packageTypeHandler := packageTypeAdapter.NewPackageTypeHandler(packageTypeUsecase)

	packageType := app.Group("/package-types")
	packageType.Get("/", packageTypeHandler.GetAll)

	protected := packageType.Group("/")
	protected.Use(middleware.IsAuth)

	admin := protected.Group("/")
	admin.Use(middleware.IsAdmin)
	admin.Post("/", packageTypeHandler.Insert)
	admin.Get("/:id", packageTypeHandler.GetByID)
	admin.Patch("/:id", packageTypeHandler.Update)
	admin.Delete("/:id", packageTypeHandler.Delete)
}
