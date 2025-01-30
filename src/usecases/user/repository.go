package userUsecase

import Entities "github.com/onosannnnt/bonbaan-BE/src/entities"

// ส่วนที่กำหนดการทำงานของ repository
type UserRepository interface {
	Insert(user *Entities.User) error
	FindByEmailOrUsername(user *Entities.User) (*Entities.User, error)
	FindByID(id *string) (*Entities.User, error)
	Update(user *Entities.User) (*Entities.User, error)
	FindAll() (*[]Entities.User, error)
	Delete(id *string) error
}

type OtpRepository interface {
	Insert(otp *Entities.Otp) error
	GetByEmail(email *string, otp *string) (*Entities.Otp, error)
	DeleteByEmail(email *string) error
}
