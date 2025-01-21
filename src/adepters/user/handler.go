package userAdepter

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/onosannnnt/bonbaan-BE/src/Constance"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	userUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/user"
)

// ส่วนที่ต่อกับ input handler
type UserHandler struct {
	userUsecase userUsecase.UserUsecase
}

// สร้าง instance ของ UserHandler
func NewUserHandler(userUsecase userUsecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// UserHandler function
func (h *UserHandler) Register() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := Entities.User{}
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"message": "Please fill all the require fields",
				"error":   err.Error(),
			})
		}
		if err := h.userUsecase.Register(user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "success",
		})
	}
}

// Custom Body Request
type LoginRequest struct {
	EmailOrUsername string `json:"emailOrUsername"`
	Password        string `json:"password"`
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var loginRequest LoginRequest
	var user Entities.User
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"message": "Please fill all the require fields",
			"error":   err.Error(),
		})
	}
	user.Email = loginRequest.EmailOrUsername
	user.Username = loginRequest.EmailOrUsername
	user.Password = loginRequest.Password
	token, err := h.userUsecase.Login(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err.Error(),
		})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * 3),
		Secure:   true,
		SameSite: "None",
		HTTPOnly: true,
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

func (h *UserHandler) Me(c *fiber.Ctx) error {
	user, err := h.userUsecase.Me(c.Locals(Constance.UserID_ctx).(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"user":    user,
	})
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	// Attempt to clear the cookie
	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Expires: time.Now().Add(-time.Hour * 24),
		Value:   "",
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
