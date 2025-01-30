package userUsecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/onosannnnt/bonbaan-BE/src/Config"
	"github.com/onosannnnt/bonbaan-BE/src/Constance"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

// ส่วนที่ต่อกับ input handler
type UserUsecase interface {
	InsertOTP(user *Entities.User) error
	Register(user *model.VerifyUserRequest) error
	Login(user *Entities.User) (token *string, err error)
	Me(userId *string) (user *Entities.User, err error)
	ChangePassword(userId *string, password *model.ChangePasswordRequest) (*Entities.User, error)
	GetAll() (*[]Entities.User, error)
	GetByID(userId *string) (*Entities.User, error)
	GetByEmailOrUsername(user *Entities.User) (*Entities.User, error)
	Delete(userId *string) error
	Update(user *model.UpdateRequest) (*Entities.User, error)
}

// ส่วนที่ต่อกับ driver handler
type UserService struct {
	userRepo UserRepository
	otpRepo  OtpRepository
}

// สร้าง instance ของ UserService
func NewUserService(userRepo UserRepository, otpRepo OtpRepository) UserUsecase {
	return &UserService{
		userRepo: userRepo,
		otpRepo:  otpRepo,
	}
}

func (s *UserService) InsertOTP(user *Entities.User) error {
	s.otpRepo.DeleteByEmail(&user.Email)
	otp := &Entities.Otp{
		Email: user.Email,
		Otp: func() string {
			otp, err := utils.GenerateOTP(6)
			if err != nil {
				// handle error appropriately, for now just panic
				panic(err)
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
func (s *UserService) Register(user *model.VerifyUserRequest) error {
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
	verifyUser := &Entities.User{
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Password:  string(hashPassword),
		RoleID:    user.RoleID,
	}
	return s.userRepo.Insert(verifyUser)
}

func (s *UserService) Login(user *Entities.User) (*string, error) {
	selectUser, err := s.userRepo.FindByEmailOrUsername(user)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(selectUser.Password), []byte(user.Password)); err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{
		Constance.Email_ctx:    selectUser.Email,
		Constance.Username_ctx: selectUser.Username,
		Constance.UserID_ctx:   selectUser.ID,
		Constance.Role_ctx:     selectUser.Role.Role,
		"exp":                  time.Now().Add(time.Hour * 24 * 3).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(Config.JwtSecret))
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func (s *UserService) Me(userId *string) (*Entities.User, error) {
	selectUser, err := s.userRepo.FindByID(userId)
	if err != nil {
		return &Entities.User{}, err
	}
	return selectUser, nil
}

func (s *UserService) ChangePassword(userId *string, password *model.ChangePasswordRequest) (*Entities.User, error) {
	selectUser, err := s.userRepo.FindByID(userId)
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

func (s *UserService) Update(user *model.UpdateRequest) (*Entities.User, error) {
	selectUser, err := s.userRepo.FindByID(&user.ID)
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
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetByID(userId *string) (*Entities.User, error) {
	user, err := s.userRepo.FindByID(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetByEmailOrUsername(user *Entities.User) (*Entities.User, error) {
	selectUser, err := s.userRepo.FindByEmailOrUsername(user)
	if err != nil {
		return nil, err
	}
	return selectUser, nil
}

func (s *UserService) Delete(userId *string) error {
	return s.userRepo.Delete(userId)
}
