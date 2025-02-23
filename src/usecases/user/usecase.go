package userUsecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/onosannnnt/bonbaan-BE/src/config"
	"github.com/onosannnnt/bonbaan-BE/src/constance"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
	roleUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/role"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

// ส่วนที่ต่อกับ input handler
type UserUsecase interface {
	InsertOTP(user *Entities.User) error
	Register(user *model.CreateUserRequest) error
	Login(user *Entities.User) (token *string, err error)
	Me(UserID *string) (user *Entities.User, err error)
	ChangePassword(UserID *string, password *model.ChangePasswordRequest) (*Entities.User, error)
	SendResetPasswordMail(user *Entities.User) error
	ResetPassword(password *model.ResetPasswordRequest) (*Entities.User, error)
	GetAll() (*[]Entities.User, error)
	GetByID(UserID *string) (*Entities.User, error)
	GetByEmailOrUsername(user *Entities.User) (*Entities.User, error)
	Delete(UserID *string) error
	Update(user *model.UpdateRequest) (*Entities.User, error)
	AdminRegister(user *model.CreateUserRequest) error
}

// ส่วนที่ต่อกับ driver handler
type UserService struct {
	userRepo          UserRepository
	otpRepo           OtpRepository
	resetPasswordRepo ResetPasswordRepository
	roleRepo          roleUsecase.RoleRepository
}

// สร้าง instance ของ UserService
func NewUserService(userRepo UserRepository, otpRepo OtpRepository, resetPasswordRepo ResetPasswordRepository, roleRepo roleUsecase.RoleRepository) UserUsecase {
	return &UserService{
		userRepo:          userRepo,
		otpRepo:           otpRepo,
		resetPasswordRepo: resetPasswordRepo,
		roleRepo:          roleRepo,
	}
}

func (s *UserService) InsertOTP(user *Entities.User) error {
	s.otpRepo.DeleteByEmail(&user.Email)
	otp := &Entities.Otp{
		Email: user.Email,
		Otp: func() string {
			otp, err := utils.GenerateOTP(6)
			if err != nil {
				return err.Error()
			}
			return otp
		}(),
		Expired: time.Now().Add(time.Minute * 5),
	}
	text := fmt.Sprintf("<body><h1>Here is your OTP <b>%s<b></h1></body>", otp.Otp)
	m := gomail.NewMessage()
	m.SetHeader("From", "bonbaanofficial@gmail.com")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "One-time password!")
	m.SetBody("text/html", text)
	utils.SendingMail(m)
	return s.otpRepo.Insert(otp)
}

// ส่วนของการทำงานของ UserService
func (s *UserService) Register(user *model.CreateUserRequest) error {
	otpEntity, err := s.otpRepo.GetByEmail(&user.Email, &user.Code)
	if err != nil {
		return err
	}
	if time.Now().After(otpEntity.Expired) {
		return errors.New("otp is expired")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	role, err := s.roleRepo.GetByName(&constance.User_Role_ctx)
	if err != nil {
		return err
	}
	verifyUser := &Entities.User{
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Password:  string(hashPassword),
		Phone:     user.Phone,
		RoleID:    role.ID,
	}
	return s.userRepo.Insert(verifyUser)
}

func (s *UserService) Login(user *Entities.User) (*string, error) {
	selectUser, err := s.userRepo.GetByEmailOrUsername(user)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(selectUser.Password), []byte(user.Password)); err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{
		constance.Email_ctx:    selectUser.Email,
		constance.Username_ctx: selectUser.Username,
		constance.UserID_ctx:   selectUser.ID,
		constance.Role_ctx:     selectUser.Role.Role,
		"exp":                  time.Now().Add(time.Hour * 24 * 3).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JwtSecret))
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func (s *UserService) Me(UserID *string) (*Entities.User, error) {
	selectUser, err := s.userRepo.GetByID(UserID)
	if err != nil {
		return &Entities.User{}, err
	}
	return selectUser, nil
}

func (s *UserService) ChangePassword(UserID *string, password *model.ChangePasswordRequest) (*Entities.User, error) {
	selectUser, err := s.userRepo.GetByID(UserID)
	if err != nil {
		return &Entities.User{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(selectUser.Password), []byte(password.OldPassword)); err != nil {
		return &Entities.User{}, err
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return &Entities.User{}, err
	}
	selectUser.Password = string(hashPassword)
	selectUser, err = s.userRepo.Update(selectUser)
	if err != nil {
		return &Entities.User{}, err
	}
	return selectUser, nil
}

func (s *UserService) SendResetPasswordMail(user *Entities.User) error {
	selectUser, err := s.userRepo.GetByEmailOrUsername(user)
	idStr := selectUser.ID.String()
	s.resetPasswordRepo.DeleteByID(&idStr)
	if err != nil {
		return err
	}
	resetPassword := &Entities.ResetPassword{
		UserID: selectUser.ID,
		Code: func() string {
			token, err := utils.GenerateOTP(6)
			if err != nil {
				return ""
			}
			return token
		}(),
		Expired: time.Now().Add(time.Minute * 5),
	}
	text := fmt.Sprintf("<body><h1>Here is your password reset Code <b>%s<b></h1></body>", resetPassword.Code)
	m := gomail.NewMessage()
	m.SetHeader("From", "bonbaanofficial@gmail.com")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "One-time password!")
	m.SetBody("text/html", text)
	utils.SendingMail(m)
	return s.resetPasswordRepo.Insert(resetPassword)
}

func (s *UserService) ResetPassword(password *model.ResetPasswordRequest) (*Entities.User, error) {
	selectUser, err := s.userRepo.GetByEmailOrUsername(&Entities.User{Email: password.Email})
	if err != nil {
		return nil, err
	}
	idStr := selectUser.ID.String()
	resetPassword, err := s.resetPasswordRepo.GetByID(&idStr, &password.Code)
	if err != nil {
		return nil, err
	}
	fmt.Println(resetPassword.User.Email)
	for password.Email != resetPassword.User.Email {
		idStr := selectUser.ID.String()
		resetPassword, err = s.resetPasswordRepo.GetByID(&idStr, &password.Code)
		if err != nil {
			return nil, err
		}
	}
	if time.Now().After(resetPassword.Expired) {
		return nil, errors.New("reset password is expired")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	selectUser.Password = string(hashPassword)
	s.resetPasswordRepo.DeleteByID(&idStr)
	selectUser, err = s.userRepo.Update(selectUser)
	if err != nil {
		return nil, err
	}
	text := "<body><h1>Your password has been changed</h1></body>"
	m := gomail.NewMessage()
	m.SetHeader("From", "bonbaanofficial@gmail.com")
	m.SetHeader("To", selectUser.Email)
	m.SetHeader("Subject", "One-time password!")
	m.SetBody("text/html", text)
	utils.SendingMail(m)
	return selectUser, nil
}

func (s *UserService) Update(user *model.UpdateRequest) (*Entities.User, error) {
	selectUser, err := s.userRepo.GetByID(&user.ID)
	if err != nil {
		return nil, err
	}
	selectUser.Username = user.Username
	selectUser.Firstname = user.FirstName
	selectUser.Lastname = user.LastName
	selectUser.Email = user.Email
	selectUser.Role.Role = user.Role
	selectUser, err = s.userRepo.Update(selectUser)
	if err != nil {
		return nil, err
	}
	return selectUser, nil
}
func (s *UserService) GetAll() (*[]Entities.User, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetByID(UserID *string) (*Entities.User, error) {
	user, err := s.userRepo.GetByID(UserID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetByEmailOrUsername(user *Entities.User) (*Entities.User, error) {
	selectUser, err := s.userRepo.GetByEmailOrUsername(user)
	if err != nil {
		return nil, err
	}
	return selectUser, nil
}

func (s *UserService) Delete(UserID *string) error {
	return s.userRepo.Delete(UserID)
}

func (s *UserService) AdminRegister(user *model.CreateUserRequest) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	role, err := s.roleRepo.GetByName(&constance.Admin_Role_ctx)
	if err != nil {
		return err
	}
	verifyUser := &Entities.User{
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Password:  string(hashPassword),
		Phone:     user.Phone,
		RoleID:    role.ID,
	}
	return s.userRepo.Insert(verifyUser)
}
