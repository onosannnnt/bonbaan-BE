package router

import (
	"github.com/gofiber/fiber/v2"
	packageAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/package"
	packageUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/package"
	"gorm.io/gorm"
)

func InitPackageRouter(app *fiber.App, db *gorm.DB) {

	packageRepo := packageAdapter.NewPackageDriver(db)
	packageUsecase := packageUsecase.NewPackageUsecase(packageRepo)
	packageHandler := packageAdapter.NewPackageHandler(packageUsecase)

	packages := app.Group("/packages")
	packages.Get("/", packageHandler.GetAllPackage)
	packages.Get("/:id", packageHandler.GetByPackageID)

	// protected := packages.Group("/")
	// protected.Use(middleware.IsAuth)

	// admin := protected.Group("/")
	// admin.Use(middleware.IsAdmin)



	packages.Post("/", packageHandler.CreatePackage)

	packages.Patch("/:id", packageHandler.UpdatePackage)
	
	packages.Delete("/:id", packageHandler.DeletePackage)

	// package.Get("/service/:serviceID", packageHandler.GetByServiceID) 




}