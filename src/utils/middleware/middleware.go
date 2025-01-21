package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/onosannnnt/bonbaan-BE/src/Config"
	"github.com/onosannnnt/bonbaan-BE/src/Constance"
)

func IsAuth(c *fiber.Ctx) error {
	cookie := c.Cookies("token")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	token, err := jwt.ParseWithClaims(cookie, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Config.JwtSecret), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
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
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	return c.Next()
}

type Owner struct {
	UserId string `json:"owner"`
}

// on implementation. I did sure this should be work. if it work, it work. if is not, read this again
func IsOwner(c *fiber.Ctx) error {
	var owner Owner
	if err := c.BodyParser(&owner); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"error":   err.Error(),
		})
	}
	userID, ok := c.Locals(Constance.UserID_ctx).(string)
	if !ok && userID != owner.UserId {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	return c.Next()
}
