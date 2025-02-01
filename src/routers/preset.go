package router

import (
	"github.com/gofiber/fiber/v2"
	presetAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/preset"
	presetUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/preset"
	"github.com/onosannnnt/bonbaan-BE/src/utils/middleware"
	"gorm.io/gorm"
)

func InitPresetRouter(app *fiber.App, db *gorm.DB) {

	presetRepo := presetAdapter.NewPresetDriver(db)
	presetUsecase := presetUsecase.NewPresetUsecase(presetRepo)
	presetHandler := presetAdapter.NewPresetHandler(presetUsecase)

	preset := app.Group("/preset")
	preset.Get("/", presetHandler.GetAllPreset)
	preset.Post("/", presetHandler.CreatePreset)
	preset.Get("/:id", presetHandler.GetByPresetID)
	preset.Patch("/:id", presetHandler.UpdatePreset)
	preset.Delete("/:id", presetHandler.DeletePreset)
	preset.Get("/service/:serviceID", presetHandler.GetByServiceID) // Add this line to handle the GetByServiceID endpoint


	protected := preset.Group("/protected")
	protected.Use(middleware.IsAuth)

	admin := protected.Group("/admin")
	admin.Use(middleware.IsAdmin)


}