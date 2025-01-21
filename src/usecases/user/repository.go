package userUsecase

import Entities "github.com/onosannnnt/bonbaan-BE/src/entities"

// ส่วนที่กำหนดการทำงานของ repository
type UserRepository interface {
	Insert(user Entities.User) error
	FindByEmailOrUsername(user *Entities.User) (*Entities.User, error)
	FindByID(id *string) (*Entities.User, error)
}
