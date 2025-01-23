package Entities

import (
	"gorm.io/gorm"
)

func InitEntity(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}, &Role{}, &Order{}, &Status{}, &Order{}); err != nil {
		panic(err)
	}
	return nil
}
