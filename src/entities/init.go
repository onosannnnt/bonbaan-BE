package Entities

import (
	"gorm.io/gorm"
)

func InitEntity(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}, &Role{}, &Order{}, &Service{}, &Status{}, &Order{}, &Package{}, &Otp{}, &ResetPassword{}, &Category{}, &Attachment{}, &Transaction{}); err != nil {
		panic(err)
	}
	return nil
}
