package userUsecase

import Entities "github.com/onosannnnt/bonbaan-BE/src/entities"

// ส่วนที่กำหนดการทำงานของ repository
type UserRepository interface {
	Insert(user *Entities.User) error
    InsertInterest(interests *[]Entities.Interest, userID *string) error  // <-- changed parameter type
	GetByEmailOrUsername(user *Entities.User) (*Entities.User, error)
	GetByID(id *string) (*Entities.User, error)
	Update(user *Entities.User) (*Entities.User, error)
	GetAll() (*[]Entities.User, error)
	Delete(id *string) error
	GetInterestByUserID(id *string) (*Entities.User, error)
	DeleteInterest(userID *string, categoryID *string) error
}

type OtpRepository interface {
	Insert(otp *Entities.Otp) error
	GetByEmail(id *string, otp *string) (*Entities.Otp, error)
	DeleteByEmail(email *string) error
}

type ResetPasswordRepository interface {
	Insert(resetPassword *Entities.ResetPassword) error
	GetByID(id *string, token *string) (*Entities.ResetPassword, error)
	DeleteByID(id *string) error
}
