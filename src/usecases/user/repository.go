package userUsecase

import Entities "github.com/onosannnnt/bonbaan-BE/src/entities"

// ส่วนที่กำหนดการทำงานของ repository
type UserRepository interface {
	Insert(user *Entities.User) error
	GetByEmailOrUsername(user *Entities.User) (*Entities.User, error)
	GetByID(id *string) (*Entities.User, error)
	Update(user *Entities.User) (*Entities.User, error)
	GetAll() (*[]Entities.User, error)
	Delete(id *string) error
}

type OtpRepository interface {
	Insert(otp *Entities.Otp) error
	GetByEmail(email *string, otp *string) (*Entities.Otp, error)
	DeleteByEmail(email *string) error
}

type ResetPasswordRepository interface {
	Insert(resetPassword *Entities.ResetPassword) error
	GetByToken(token *string) (*Entities.ResetPassword, error)
	DeleteByEmail(email *string) error
}
