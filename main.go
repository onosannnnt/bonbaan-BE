package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/onosannnnt/bonbaan-BE/src/Config"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	router "github.com/onosannnnt/bonbaan-BE/src/routers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	// Connect to database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", Config.DbHost, Config.DbPort, Config.DbUser, Config.DbPassword, Config.DbSchema)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// Postgres install uuid-ossp extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		panic("failed to create uuid-ossp extension")
	}

	//Initialize Entities
	Entities.InitEntity(db)

	if err != nil {
		panic("failed to connect database")
	}

	app := fiber.New()

	router.InitUserRouter(app, db)
	router.InitRoleRouter(app, db)

	app.Listen(":" + Config.Port)
}

func Hello(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Hello, World!"})
}
