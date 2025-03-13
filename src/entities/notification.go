package Entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:(uuid_generate_v4())"`
	Header  string    `json:"header"`
	Body    string    `json:"body"`
	IsRead  bool      `json:"is_read"`
	OrderID uuid.UUID `json:"orderID" gorm:"type:uuid;default:null"`
	Order   Order     `gorm:"foreignKey:OrderID;default:null"`
	UserID  uuid.UUID `json:"userID" gorm:"type:uuid"`
	User    User      `gorm:"foreignKey:UserID "`
}
