package router_test

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	router "github.com/onosannnnt/bonbaan-BE/src/routers"
)

func TestInitUserRouter(t *testing.T) {
    // Setup: create a new Fiber app and a dummy GORM DB instance.
    app := fiber.New()
    db := &gorm.DB{} // Using a dummy DB as the router only uses it for initialization.

    // Execute: initialize the user router.
    router.InitUserRouter(app, db)

    // Get the route stack; note that Stack() returns [][]*Route.
    stack := app.Stack()

    // Build a map where keys are route paths and values are slices 
    // of registered HTTP methods. Use route.Method as the method.
    routeMap := make(map[string][]string)
    for _, routes := range stack {
        for _, route := range routes {
            routeMap[route.Path] = append(routeMap[route.Path], route.Method)
        }
    }

    // Define expected routes with path as key and a slice of expected HTTP methods.
    expectedRoutes := map[string][]string{
        "/users/send-otp":                      {"POST"},
        "/users/register":                      {"POST"},
        "/users/login":                         {"POST"},
        "/users/send-reset-password":           {"POST"},
        "/users/reset-password/":               {"POST"},
        "/users/me":                            {"GET"},
        "/users/":                              {"GET", "DELETE", "PATCH"}, // Note: PATCH is registered twice in your router.
        "/users/:id":                           {"GET"},
        "/users/email-or-username/:emailOrUsername": {"GET"},
    }

    // Check that each expected method is registered for each expected route.
    for path, methods := range expectedRoutes {
        registeredMethods, exists := routeMap[path]
        assert.True(t, exists, "Route %s not found", path)
        for _, m := range methods {
            found := false
            for _, rm := range registeredMethods {
                if rm == m {
                    found = true
                    break
                }
            }
            assert.True(t, found, "Expected method %s for route %s not found", m, path)
        }
    }
}