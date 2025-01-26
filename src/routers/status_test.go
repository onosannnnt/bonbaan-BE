package router


import (
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestInitStatusRouter(t *testing.T) {
	// Setup
	app := fiber.New()
	db := &gorm.DB{}

	// Test execution
	InitStatusRouter(app, db)

	// Verify routes are registered
	routes := app.GetRoutes()
	
	expectedRoutes := map[string]string{
		"/status/":              "GET",
		"/status/:id":          "GET", 
		"/status/name/:name":   "GET",
		"/status/protected/admin/": "POST",
		"/status/protected/admin/:id": "PUT,DELETE",
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
