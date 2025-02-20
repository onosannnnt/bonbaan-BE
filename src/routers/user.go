package router

import (
	// jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	// "github.com/onosannnnt/bonbaan-BE/src/Config"
	otpDriver "github.com/onosannnnt/bonbaan-BE/src/adepters/otp"
	resetpasswordDriver "github.com/onosannnnt/bonbaan-BE/src/adepters/reset_password"
	userAdepter "github.com/onosannnnt/bonbaan-BE/src/adepters/user"
	userUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/user"
	"github.com/onosannnnt/bonbaan-BE/src/utils/middleware"
	"gorm.io/gorm"
)

func InitUserRouter(app *fiber.App, db *gorm.DB) {

	userRepo := userAdepter.NewUserDriver(db)
	otpRepo := otpDriver.NewOtpDriver(db)
	resetPasswordRepo := resetpasswordDriver.NewOtpDriver(db)
	userUsecase := userUsecase.NewUserService(userRepo, otpRepo, resetPasswordRepo)
	userHandler := userAdepter.NewUserHandler(userUsecase)

	user := app.Group("/users")
	user.Post("/send-otp", userHandler.SendOTP)
	user.Post("/register", userHandler.Register())
	user.Post("/login", userHandler.Login)
	user.Post("/send-reset-password", userHandler.SendResetPasswordMail)
	user.Post("/reset-password/", userHandler.ResetPassword)

	protect := user.Group("/")
	protect.Use(middleware.IsAuth)
	
	protect.Get("/me", userHandler.Me)
	protect.Get("/", userHandler.GetAll)
	protect.Get("/:id", userHandler.GetByID)
	protect.Get("/email-or-username/:emailOrUsername", userHandler.GetByEmailOrUsername)
	protect.Delete("/", userHandler.Delete)

	owner := protect.Group("/")
	owner.Use(middleware.IsOwner)

	owner.Patch("/", userHandler.ChangePassword)
	owner.Patch("/", userHandler.Update)

}
