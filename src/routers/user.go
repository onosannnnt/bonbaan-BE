package router

import (
	// jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	// "github.com/onosannnnt/bonbaan-BE/src/Config"
	otpDriver "github.com/onosannnnt/bonbaan-BE/src/adepters/otp"
	userAdepter "github.com/onosannnnt/bonbaan-BE/src/adepters/user"
	userUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/user"
	"github.com/onosannnnt/bonbaan-BE/src/utils/middleware"
	"gorm.io/gorm"
)

func InitUserRouter(app *fiber.App, db *gorm.DB) {

	userRepo := userAdepter.NewUserDriver(db)
	otpRepo := otpDriver.NewOtpDriver(db)
	userUsecase := userUsecase.NewUserService(userRepo, otpRepo)
	userHandler := userAdepter.NewUserHandler(userUsecase)

	user := app.Group("/users")
	user.Post("/send-otp", userHandler.SendOTP)
	user.Post("/register", userHandler.Register())
	user.Post("/login", userHandler.Login)

	protect := user.Group("/protected")
	protect.Use(middleware.IsAuth)
	protect.Get("/me", userHandler.Me)
	protect.Get("/", userHandler.GetAll)
	protect.Get("/:id", userHandler.GetByID)
	protect.Get("/email-or-username/:emailOrUsername", userHandler.GetByEmailOrUsername)
	protect.Delete("/", userHandler.Delete)

	owner := protect.Group("/owner")
	owner.Use(middleware.IsOwner)

	owner.Put("/change-password", userHandler.ChangePassword)
	owner.Put("/", userHandler.Update)

}
