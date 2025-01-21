package router

import (
	"github.com/gofiber/fiber/v2"
	userAdepter "github.com/onosannnnt/bonbaan-BE/src/adepters/user"
	userUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/user"
	"github.com/onosannnnt/bonbaan-BE/src/utils/middleware"
	"gorm.io/gorm"
)

func InitUserRouter(app *fiber.App, db *gorm.DB) {

	userRepo := userAdepter.NewUserDriver(db)
	userUsecase := userUsecase.NewUserService(userRepo)
	userHandler := userAdepter.NewUserHandler(userUsecase)

	user := app.Group("/users")
	user.Post("/register", userHandler.Register())
	user.Post("/login", userHandler.Login)

	protect := user.Group("/protected")
	protect.Use(middleware.IsAuth)
	protect.Get("/me", userHandler.Me)
	protect.Get("/logout", userHandler.Logout)
}
