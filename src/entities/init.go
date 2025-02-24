package Entities

import (
	"gorm.io/gorm"
)

func InitEntity(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}, &Role{}, &Order{}, &Service{}, &Status{}, &Order{}, &Package{}, &Otp{}, &ResetPassword{}, &Category{}, &Attachment{}, &Transaction{},,&PackageType{}); err != nil {
		panic(err)
	}
	if err := InitializeRoleData(db); err != nil {
		panic(err)
	}
	if err := InitializeStatusData(db); err != nil {
		panic(err)
	}
	if err := InitializeUserData(db); err != nil {
		panic(err)
	}
	return nil
}
