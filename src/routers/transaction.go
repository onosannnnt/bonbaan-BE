package router

import (
	"github.com/gofiber/fiber/v2"
	transactionAdepter "github.com/onosannnnt/bonbaan-BE/src/adepters/transaction"
	transactionUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/transaction"
	"github.com/onosannnnt/bonbaan-BE/src/utils/middleware"
	"gorm.io/gorm"
)

func InitTransactionRouter(app *fiber.App, db *gorm.DB) {
	transactionRepo := transactionAdepter.NewTransactionDriver(db)
	transactionUsecase := transactionUsecase.NewTransectionService(transactionRepo)
	transactionHandler := transactionAdepter.NewStatusHandler(transactionUsecase)

	transaction := app.Group("/transaction")

	protected := transaction.Group("/")
	protected.Use(middleware.IsAuth)
	protected.Get("/", transactionHandler.GetAllTransaction)
	protected.Get("/:id", transactionHandler.GetTransactionByID)
	protected.Post("/", transactionHandler.InsertTransaction)
	protected.Patch("/:id", transactionHandler.UpdateTransaction)
	protected.Delete("/:id", transactionHandler.DeleteTransaction)

}
