package Entities

import (
	"github.com/google/uuid"
	"github.com/onosannnnt/bonbaan-BE/src/constance"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;primaryKey;default:(uuid_generate_v4())"`
	Role string    `json:"role"`
}

func InitializeData(db *gorm.DB) error {
	roles := []Role{
		{Role: constance.Admin_Role_ctx},
		{Role: "user"},
	}
	for _, role := range roles {
		if err := db.FirstOrCreate(&role).Error; err != nil {
			return err
		}
	}
	return nil
}
