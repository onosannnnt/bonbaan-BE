package userUsecase

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/onosannnnt/bonbaan-BE/src/Config"
	"github.com/onosannnnt/bonbaan-BE/src/Constance"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
	"golang.org/x/crypto/bcrypt"
)

// ส่วนที่ต่อกับ input handler
type UserUsecase interface {
	Register(user Entities.User) error
	Login(user Entities.User) (token string, err error)
	Logout(token string) error
	Me(userId string) (user Entities.User, err error)
	ChangePassword(userId string, password model.ChangePasswordRequest) (user Entities.User, err error)
	GetAll() ([]Entities.User, error)
	Delete(userId string) error
}

// ส่วนที่ต่อกับ driver handler
type UserService struct {
	userRepo UserRepository
}

// สร้าง instance ของ UserService
func NewUserService(repo UserRepository) UserUsecase {
	return &UserService{
		userRepo: repo,
	}
}

// ส่วนของการทำงานของ UserService
func (s *UserService) Register(user Entities.User) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashPassword)
	return s.userRepo.Insert(user)
}

func (s *UserService) Login(user Entities.User) (string, error) {
	selectUser, err := s.userRepo.FindByEmailOrUsername(&user)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(selectUser.Password), []byte(user.Password)); err != nil {
		return "", err
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims[Constance.Email_ctx] = selectUser.Email
	claims[Constance.Username_ctx] = selectUser.Username
	claims[Constance.UserID_ctx] = selectUser.ID
	claims[Constance.Role_ctx] = selectUser.Role.Role
	claims["exp"] = time.Now().Add(time.Hour * 24 * 3).Unix()
	tokenString, err := token.SignedString([]byte(Config.JwtSecret))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *UserService) Logout(token string) error {
	return nil
}

func (s *UserService) Me(userId string) (Entities.User, error) {
	selectUser, err := s.userRepo.FindByID(&userId)
	if err != nil {
		return Entities.User{}, err
	}
	return *selectUser, nil
}

func (s *UserService) ChangePassword(userId string, password model.ChangePasswordRequest) (Entities.User, error) {
	selectUser, err := s.userRepo.FindByID(&userId)
	if err != nil {
		return Entities.User{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(selectUser.Password), []byte(password.OldPassword)); err != nil {
		return Entities.User{}, err
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return Entities.User{}, err
	}
	selectUser.Password = string(hashPassword)
	selectUser, err = s.userRepo.Update(*selectUser)
	if err != nil {
		return Entities.User{}, err
	}
	return *selectUser, nil
}

func (s *UserService) GetAll() ([]Entities.User, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (s *UserService) Delete(userId string) error {
	return s.userRepo.Delete(&userId)
}
