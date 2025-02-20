package router_test

import (
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	router "github.com/onosannnnt/bonbaan-BE/src/routers"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestInitOrderRouter(t *testing.T) {
    // Setup
    app := fiber.New()
    db := &gorm.DB{}

    // Test initialization
    router.InitOrderRouter(app, db)

    // Verify routes are registered
    routes := app.GetRoutes()
    
    expectedRoutes := map[string]string{
        "/orders/":     "POST,GET",
        "/orders/:id": "GET,PATCH,DELETE",
    }

    // Assert routes exist
    for path, methods := range expectedRoutes {
        for _, method := range strings.Split(methods, ",") {
            found := false
            for _, route := range routes {
                if route.Path == path && route.Method == method {
                    found = true
                    break
                }
            }
            assert.True(t, found, "Expected route %s %s not found", method, path)
        }
    }
}

