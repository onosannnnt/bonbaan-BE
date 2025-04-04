package router

import (
	"github.com/gofiber/fiber/v2"
	orderTypeAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/order_type"
	orderTypeUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/order_type"
	"github.com/onosannnnt/bonbaan-BE/src/utils/middleware"
	"gorm.io/gorm"
)

func InitOrderTypeRouter(app *fiber.App, db *gorm.DB) {

	orderTypeRepo := orderTypeAdapter.NewOrderTypeDriver(db)
	orderTypeUsecase := orderTypeUsecase.NewOrderTypeService(orderTypeRepo)
	orderTypeHandler := orderTypeAdapter.NewOrderTypeHandler(orderTypeUsecase)

	orderType := app.Group("/order-types")
	orderType.Get("/", orderTypeHandler.GetAll)

	protected := orderType.Group("/")
	protected.Use(middleware.IsAuth)

	admin := protected.Group("/")
	admin.Use(middleware.IsAdmin)
	admin.Post("/", orderTypeHandler.Insert)
	admin.Get("/:id", orderTypeHandler.GetByID)
	admin.Patch("/:id", orderTypeHandler.Update)
	admin.Delete("/:id", orderTypeHandler.Delete)
}
