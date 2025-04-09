package Entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Interest struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:(uuid_generate_v4())"`
	UserID    uuid.UUID `json:"userID"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
	CategoryID uuid.UUID `json:"categoryID"`
	Category  Category  `gorm:"foreignKey:CategoryID;references:ID"`
}