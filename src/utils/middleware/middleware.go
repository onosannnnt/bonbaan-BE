package middleware

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/onosannnnt/bonbaan-BE/src/Config"
	"github.com/onosannnnt/bonbaan-BE/src/Constance"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
)

func IsAuth(c *fiber.Ctx) error {
	jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(Config.JwtSecret)},
		ContextKey: "jwt",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return utils.ResponseJSON(c, fiber.StatusForbidden, "Unauthorized", err, nil)
		},
	})
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		authHeader = authHeader[len("Bearer "):]
	}
	token, err := jwt.ParseWithClaims(authHeader, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Config.JwtSecret), nil
	})

	if err != nil || !token.Valid {
		return utils.ResponseJSON(c, fiber.StatusUnauthorized, "Unauthorized", err, nil)
	}
	claim := token.Claims.(jwt.MapClaims)
	c.Locals(Constance.UserID_ctx, claim[Constance.UserID_ctx])
	c.Locals(Constance.Email_ctx, claim[Constance.Email_ctx])
	c.Locals(Constance.Username_ctx, claim[Constance.Username_ctx])
	c.Locals(Constance.Role_ctx, claim[Constance.Role_ctx])
	return c.Next()
}

func IsAdmin(c *fiber.Ctx) error {
	role, ok := c.Locals(Constance.Role_ctx).(string)

	if !ok && role != Constance.Admin_Role_ctx {
		return utils.ResponseJSON(c, fiber.StatusForbidden, "Forbidden", nil, nil)
	}
	return c.Next()
}

// on implementation. I did sure this should be work. if it work, it work. if is not, read this again
func IsOwner(c *fiber.Ctx) error {
	owner := c.Get("UserID")
	userID, ok := c.Locals(Constance.UserID_ctx).(string)
	if !ok || userID != owner {
		return utils.ResponseJSON(c, fiber.StatusForbidden, "Forbidden", nil, nil)
	}
	return c.Next()
}
