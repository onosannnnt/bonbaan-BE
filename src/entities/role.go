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

func InitializeRoleData(db *gorm.DB) error {
	roles := []Role{
		{Role: constance.User_Role_ctx},
		{Role: constance.Admin_Role_ctx},
	}
	for _, role := range roles {
		if err := db.FirstOrCreate(&role, Role{Role: role.Role}).Error; err != nil {
			return err
		}
	}
	return nil
}
