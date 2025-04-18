package Entities

import (
	"gorm.io/gorm"
)

func InitEntity(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}, &Role{}, &Order{}, &Service{}, &Status{}, &Package{}, &Otp{}, &ResetPassword{}, &Category{}, &Attachment{}, &Transaction{}, &OrderType{}, &Review{}, &Review_utils{}, &Notification{}, &VowRecord{},&Recommendation{},&RecommendationUtil{},&Interest{}); err != nil {
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
	if err := InitializeOrderTypeData(db); err != nil {
		panic(err)
	}
	return nil
}
