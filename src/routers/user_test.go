package router

import (
	"strings"
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
	
	expectedRoutes := map[string]string{
		"/users/register":                      "POST",
		"/users/login":                         "POST",
		"/users/protected/me":                  "GET",
		"/users/protected":                     "GET,DELETE",
		"/users/protected/owner/change-password": "PUT",
		"/users/protected/owner":               "PUT",
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
	// Verify middleware
	stack := app.Stack()
	hasAuthMiddleware := false
	hasOwnerMiddleware := false

	for _, s := range stack {
		for _, r := range s {
			if r.Path == "/users/protected" {
				hasAuthMiddleware = true
			}
			if r.Path == "/users/protected/owner" {
				hasOwnerMiddleware = true 
			}
		}
	}

	assert.True(t, hasAuthMiddleware, "Auth middleware not found")
	assert.True(t, hasOwnerMiddleware, "Owner middleware not found")
}