package main

import (
	"fmt"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/onosannnnt/bonbaan-BE/src/config"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	router "github.com/onosannnnt/bonbaan-BE/src/routers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// image_upload()
	// Connect to database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.DbHost, config.DbPort, config.DbUser, config.DbPassword, config.DbSchema)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// Postgres install uuid-ossp extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		panic("failed to create uuid-ossp extension")
	}

	//Initialize Entities
	Entities.InitEntity(db)

	//check Entities in database

	if err != nil {
		panic("failed to connect database")
	}

	app := fiber.New(fiber.Config{
		BodyLimit: math.MaxInt64,
	})

	app.Use(cors.New(
		cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		},
	))

	router.InitUserRouter(app, db)
	router.InitRoleRouter(app, db)
	router.InitStatusRouter(app, db)
	router.InitOrderRouter(app, db)
	router.ServiceRouter(app, db)
	router.InitTransactionRouter(app, db)
	router.InitPackageRouter(app, db)
	router.InitCategoryRouter(app, db)
	router.InitAttachmentRouter(app, db)
	router.InitPackageTypeRouter(app, db)
	app.Listen(":" + config.Port)
	
}

func Hello(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Hello, World!"})
}
