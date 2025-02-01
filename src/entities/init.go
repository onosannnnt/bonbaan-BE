package Entities

import (
	"gorm.io/gorm"
)

func InitEntity(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}, &Role{}, &Order{}, &Service{}, &Status{}, &Order{},&Preset{} ,&Otp{}, &ResetPassword{}); err != nil {
		panic(err)
	}
	return nil
}
