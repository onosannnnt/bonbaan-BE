package router

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestInitUserRouter(t *testing.T) {
	// Setup
	app := fiber.New()
	db := &gorm.DB{}

	// Test execution
	InitUserRouter(app, db)

	// Verify routes are registered
	routes := app.GetRoutes()

	// Test /users routes
	assert.Contains(t, getRoutes(routes), "POST /users/register")
	assert.Contains(t, getRoutes(routes), "POST /users/login")

	// Test protected routes
	assert.Contains(t, getRoutes(routes), "GET /users/protected/me")
	assert.Contains(t, getRoutes(routes), "GET /users/protected")
	assert.Contains(t, getRoutes(routes), "DELETE /users/protected")

	// Test owner routes
	assert.Contains(t, getRoutes(routes), "PUT /users/protected/owner/change-password")
	assert.Contains(t, getRoutes(routes), "PUT /users/protected/owner")
}

func getRoutes(routes []fiber.Route) []string {
	var routeStrings []string
	for _, route := range routes {
		routeStrings = append(routeStrings, route.Method+" "+route.Path)
	}
	return routeStrings
}