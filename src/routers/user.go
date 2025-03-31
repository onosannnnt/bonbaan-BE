package router

import (
	"github.com/gofiber/fiber/v2"
	categoryAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/category"
	notificationAdepter "github.com/onosannnnt/bonbaan-BE/src/adepters/notification"
	orderAdepter "github.com/onosannnnt/bonbaan-BE/src/adepters/order"
	otpDriver "github.com/onosannnnt/bonbaan-BE/src/adepters/otp"
	resetpasswordDriver "github.com/onosannnnt/bonbaan-BE/src/adepters/reset_password"
	roleAdapter "github.com/onosannnnt/bonbaan-BE/src/adepters/role"
	userAdepter "github.com/onosannnnt/bonbaan-BE/src/adepters/user"
	notificationUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/notification"
	orderUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/order"
	userUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/user"
	"github.com/onosannnnt/bonbaan-BE/src/utils/middleware"
	"gorm.io/gorm"
)

func InitUserRouter(app *fiber.App, db *gorm.DB) {
	notificationRepo := notificationAdepter.NewNotificationDriver(db)
	notificationUsecase := notificationUsecase.NewNotificationService(notificationRepo)
	notificationHandler := notificationAdepter.NewNotificationHandler(notificationUsecase)
	orderRepo := orderAdepter.NewOrderDriver(db, nil)
	orderUsecase := orderUsecase.NewOrderService(orderRepo, nil, nil, nil, nil, nil, nil)
	orderHandler := orderAdepter.NewOrderHandler(orderUsecase)

	categoryRepo := categoryAdapter.NewCategoryDriver(db)
	userRepo := userAdepter.NewUserDriver(db)
	otpRepo := otpDriver.NewOtpDriver(db)
	roleRepo := roleAdapter.NewRoleDriver(db)
	resetPasswordRepo := resetpasswordDriver.NewOtpDriver(db)

	userUsecase := userUsecase.NewUserService(userRepo, otpRepo, resetPasswordRepo, roleRepo, categoryRepo)
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
	protect.Get("/:id/notifications", notificationHandler.GetByUserID)
	protect.Post("/:id/interest", userHandler.InsertInterest)
	protect.Get("/:id/interest", userHandler.GetInterestByUserID)
	protect.Delete("/interest/:id", userHandler.DeleteInterest)
	protect.Get("/:id/orders", orderHandler.GetByUserID)

	admin := protect.Group("/")
	admin.Use(middleware.IsAdmin)

	admin.Post("/admin-register", userHandler.AdminRegister)

}
