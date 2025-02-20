package router

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestInitServiceRouter(t *testing.T) {
    // Setup
    app := fiber.New()
    db := &gorm.DB{}

    // Test execution
    ServiceRouter(app, db)

    // Verify routes are registered
    routes := app.GetRoutes()
    
    expectedRoutes := map[string][]string{
        "/services/":    {"GET", "POST"},
        "/services/:id": {"GET", "PATCH"},
    }

    for path, expectedMethods := range expectedRoutes {
        for _, expectedMethod := range expectedMethods {
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
    hasAdminMiddleware := false

    for _, s := range stack {
        for _, r := range s {
            if r.Path == "/services/" {
                hasAuthMiddleware = true
            }
            if r.Path == "/services/:id" {
                hasAdminMiddleware = true 
            }
        }
    }

    assert.True(t, hasAuthMiddleware, "Auth middleware not found")
    assert.True(t, hasAdminMiddleware, "Admin middleware not found")
}







