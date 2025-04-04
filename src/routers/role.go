package router

import (
	"github.com/gofiber/fiber/v2"
	roleAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/role"
	roleUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/role"
	"github.com/onosannnnt/bonbaan-BE/src/utils/middleware"
	"gorm.io/gorm"
)

func InitRoleRouter(app *fiber.App, db *gorm.DB) {

	roleRepo := roleAdapter.NewRoleDriver(db)
	roleUsecase := roleUsecase.NewRoleService(roleRepo)
	roleHandler := roleAdapter.NewRoleHandler(roleUsecase)

	role := app.Group("/roles")
	role.Get("/", roleHandler.GetAllRole)

	protected := role.Group("/")
	protected.Use(middleware.IsAuth)

	admin := protected.Group("/")
	admin.Use(middleware.IsAdmin)
	admin.Post("/", roleHandler.InsertRole)
}
