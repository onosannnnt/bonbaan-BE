package Entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:(uuid_generate_v4())"`
	UserID    uuid.UUID `json:"userID"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
	ServiceID uuid.UUID `json:"serviceID"`
	Service   Service   `gorm:"foreignKey:ServiceID;references:ID"`
	OrderID   uuid.UUID `gorm:"foreignKey:OrderID;references:ID"`
	Rating    int       `json:"rating"`
	Detail    string    `json:"detail"`
}
