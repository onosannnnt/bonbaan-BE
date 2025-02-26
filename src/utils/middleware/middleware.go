package middleware

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/onosannnnt/bonbaan-BE/src/config"
	"github.com/onosannnnt/bonbaan-BE/src/constance"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
)

func IsAuth(c *fiber.Ctx) error {
	jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.JwtSecret)},
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
		return []byte(config.JwtSecret), nil
	})
	if err != nil || !token.Valid {
		return utils.ResponseJSON(c, fiber.StatusUnauthorized, "Unauthorized", err, nil)
	}
	claim := token.Claims.(jwt.MapClaims)
	c.Locals(constance.UserID_ctx, claim[constance.UserID_ctx])
	c.Locals(constance.Email_ctx, claim[constance.Email_ctx])
	c.Locals(constance.Username_ctx, claim[constance.Username_ctx])
	c.Locals(constance.Role_ctx, claim[constance.Role_ctx])
	return c.Next()
}

func IsAdmin(c *fiber.Ctx) error {
	role, ok := c.Locals(constance.Role_ctx).(string)

	if !ok || role != constance.Admin_Role_ctx {
		return utils.ResponseJSON(c, fiber.StatusForbidden, "Forbidden", nil, nil)
	}
	return c.Next()
}
