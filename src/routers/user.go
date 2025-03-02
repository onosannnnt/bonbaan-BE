package router

import (
	// jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	// "github.com/onosannnnt/bonbaan-BE/src/config"
	otpDriver "github.com/onosannnnt/bonbaan-BE/src/adepters/otp"
	resetpasswordDriver "github.com/onosannnnt/bonbaan-BE/src/adepters/reset_password"
	roleAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/role"
	userAdepter "github.com/onosannnnt/bonbaan-BE/src/adepters/user"
	userUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/user"
	"github.com/onosannnnt/bonbaan-BE/src/utils/middleware"
	"gorm.io/gorm"
)

func InitUserRouter(app *fiber.App, db *gorm.DB) {

	userRepo := userAdepter.NewUserDriver(db)
	otpRepo := otpDriver.NewOtpDriver(db)
	roleRepo := roleAdapter.NewRoleDriver(db)
	resetPasswordRepo := resetpasswordDriver.NewOtpDriver(db)
	userUsecase := userUsecase.NewUserService(userRepo, otpRepo, resetPasswordRepo, roleRepo)
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
	protect.Patch("/change-password", userHandler.ChangePassword)
	protect.Patch("/", userHandler.Update)
	protect.Post("/interest", userHandler.InsertInterest)

	admin := protect.Group("/")
	admin.Use(middleware.IsAdmin)

	admin.Post("/admin-register", userHandler.AdminRegister)

}
