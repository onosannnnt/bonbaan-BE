package userAdepter

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/onosannnnt/bonbaan-BE/src/Constance"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
	userUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/user"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
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
			return utils.ResponseJSON(c, fiber.ErrBadRequest.Code, "Please fill all the require fields", err, nil)
		}
		if err := h.userUsecase.Register(&user); err != nil {
			return utils.ResponseJSON(c, fiber.StatusConflict, "this account already exists", err, nil)
		}
		return utils.ResponseJSON(c, fiber.StatusCreated, "success", nil, nil)
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
		if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			return utils.ResponseJSON(c, fiber.ErrBadRequest.Code, "Invalid email, username or password", err, nil)
		}
		return utils.ResponseJSON(c, fiber.ErrBadRequest.Code, "Please fill all the require fields", err, nil)
	}
	user.Email = loginRequest.EmailOrUsername
	user.Username = loginRequest.EmailOrUsername
	user.Password = loginRequest.Password
	token, err := h.userUsecase.Login(&user)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusUnauthorized, "Invalid email, username or password", err, nil)
	}
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    *token,
		Expires:  time.Now().Add(time.Hour * 24 * 3),
		Secure:   true,
		SameSite: "None",
		HTTPOnly: true,
	})
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, map[string]string{"token": *token})
}

func (h *UserHandler) Me(c *fiber.Ctx) error {
	userID := c.Locals(Constance.UserID_ctx).(string)
	user, err := h.userUsecase.Me(&userID)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, user)
}

func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	changePasswordRequest := model.ChangePasswordRequest{}
	if err := c.BodyParser(&changePasswordRequest); err != nil {
		return utils.ResponseJSON(c, fiber.ErrBadRequest.Code, "Please fill all the require fields", err, nil)
	}
	userId := c.Locals(Constance.UserID_ctx).(string)
	user, err := h.userUsecase.ChangePassword(&userId, &changePasswordRequest)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, user)

}

func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	users, err := h.userUsecase.GetAll()
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, users)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	var user model.UpdateRequest
	if err := c.BodyParser(&user); err != nil {
		return utils.ResponseJSON(c, fiber.ErrBadRequest.Code, "Please fill all the require fields", err, nil)
	}

	userID, ok := c.Locals(Constance.UserID_ctx).(string)
	if !ok {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", nil, nil)
	}
	user.ID = userID
	selectUser, err := h.userUsecase.Update(&user)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, selectUser)
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals(Constance.UserID_ctx).(string)
	err := h.userUsecase.Delete(&userID)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Expires: time.Now().Add(-time.Hour * 24),
		Value:   "",
	})
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, nil)

}
