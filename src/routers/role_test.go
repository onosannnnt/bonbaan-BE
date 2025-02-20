package router

import (
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestInitRoleRouter(t *testing.T) {
	app := fiber.New()
	db := &gorm.DB{}

	InitRoleRouter(app, db)

	routes := app.GetRoutes()

	expectedRoutes := map[string]string{
		"/roles/": "GET,POST",
	}

	for path, methods := range expectedRoutes {
		for _, expectedMethod := range strings.Split(methods, ",") {
			var found bool
			for _, route := range routes {
				if route.Path == path && route.Method == expectedMethod {
					found = true
					break
				}
			}
			assert.True(t, found, "Expected route %s %s not found", expectedMethod, path)
		}
	}

}
