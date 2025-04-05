package userAdepter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/onosannnnt/bonbaan-BE/src/constance"
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

func (h *UserHandler) SendOTP(c *fiber.Ctx) error {
	user := Entities.User{}
	if err := c.BodyParser(&user); err != nil {
		return utils.ResponseJSON(c, fiber.ErrBadRequest.Code, "Please fill all the require fields", err, nil)
	}
	code, err := h.userUsecase.InsertOTP(&user)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, code)
}

// UserHandler function
func (h *UserHandler) Register() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := model.CreateUserRequest{}
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
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, map[string]string{"token": *token})
}

func (h *UserHandler) Me(c *fiber.Ctx) error {
	userID := c.Locals(constance.UserID_ctx).(string)
	if userID == "" {
		return utils.ResponseJSON(c, fiber.StatusUnauthorized, "missing UserID in header", nil, nil)
	}
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
	UserID := c.Locals(constance.UserID_ctx).(string)
	if UserID == "" {
		return utils.ResponseJSON(c, fiber.StatusUnauthorized, "missing UserID in header", nil, nil)
	}
	user, err := h.userUsecase.ChangePassword(&UserID, &changePasswordRequest)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, err.Error(), err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, user)
}

func (h *UserHandler) SendResetPasswordMail(c *fiber.Ctx) error {
	user := Entities.User{}
	if err := c.BodyParser(&user); err != nil {
		return utils.ResponseJSON(c, fiber.ErrBadRequest.Code, "Please fill all the require fields", err, nil)
	}
	if err := h.userUsecase.SendResetPasswordMail(&user); err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, nil)
}

func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	resetPasswordRequest := model.ResetPasswordRequest{}
	if err := c.BodyParser(&resetPasswordRequest); err != nil {
		return utils.ResponseJSON(c, fiber.ErrBadRequest.Code, "Please fill all the require fields", err, nil)
	}
	user, err := h.userUsecase.ResetPassword(&resetPasswordRequest)
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

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	userID := c.Params("id")
	user, err := h.userUsecase.GetByID(&userID)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, user)
}

func (h *UserHandler) GetByEmailOrUsername(c *fiber.Ctx) error {
	emailOrUsername := c.Params("emailOrUsername")
	user, err := h.userUsecase.GetByID(&emailOrUsername)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, user)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	var user model.UpdateRequest

	if err := c.BodyParser(&user); err != nil {
		return utils.ResponseJSON(c, fiber.ErrBadRequest.Code, "Please fill all the require fields", err, nil)
	}
	userID := c.Locals(constance.UserID_ctx).(string)
	if userID == "" {
		return utils.ResponseJSON(c, fiber.StatusUnauthorized, "missing UserID in header", nil, nil)
	}

	user.ID = userID
	selectUser, err := h.userUsecase.Update(&user)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, selectUser)
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals(constance.UserID_ctx).(string)
	if userID == "" {
		return utils.ResponseJSON(c, fiber.StatusUnauthorized, "missing UserID in header", nil, nil)
	}
	err := h.userUsecase.Delete(&userID)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, nil)

}

func (h *UserHandler) AdminRegister(c *fiber.Ctx) error {
	user := model.CreateUserRequest{}
	if err := c.BodyParser(&user); err != nil {
		return utils.ResponseJSON(c, fiber.ErrBadRequest.Code, "Please fill all the require fields", err, nil)
	}
	if err := h.userUsecase.AdminRegister(&user); err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, nil)
}

func (h *UserHandler) InsertInterest(c *fiber.Ctx) error {
	type InsertInterest struct {
		Interest []string `json:"categories"`
	}
	interest := InsertInterest{}
	if err := c.BodyParser(&interest); err != nil {
		return utils.ResponseJSON(c, fiber.ErrBadRequest.Code, "Please fill all the require fields", err, nil)
	}
	userID := c.Locals(constance.UserID_ctx).(string)
	if userID == "" {
		return utils.ResponseJSON(c, fiber.StatusUnauthorized, "missing UserID in header", nil, nil)
	}
	err := h.userUsecase.InsertInterest(&userID, &interest.Interest)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusCreated, "success", nil, nil)
}


func (h *UserHandler) GetInterestByUserID(c *fiber.Ctx) error {
	userID := c.Locals(constance.UserID_ctx).(string)
	if userID == "" {
		return utils.ResponseJSON(c, fiber.StatusUnauthorized, "missing UserID in header", nil, nil)
	}
	user, err := h.userUsecase.GetInterestByUserID(&userID)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, user)
}

func (h *UserHandler) DeleteInterest(c *fiber.Ctx) error {
	interest := c.Params("id")
	if interest == "" {
		return utils.ResponseJSON(c, fiber.ErrBadRequest.Code, "Please fill all the require fields", nil, nil)
	}
	userID := c.Locals(constance.UserID_ctx).(string)
	if userID == "" {
		return utils.ResponseJSON(c, fiber.StatusUnauthorized, "missing UserID in header", nil, nil)
	}
	err := h.userUsecase.DeleteInterest(&userID, &interest)
	if err != nil {
		return utils.ResponseJSON(c, fiber.StatusInternalServerError, "Internal Server Error", err, nil)
	}
	return utils.ResponseJSON(c, fiber.StatusOK, "success", nil, nil)
}
