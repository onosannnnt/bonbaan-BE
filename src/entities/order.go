package Entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JSONB map[string]interface{}

type Order struct {
	gorm.Model
	ID                 uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CancellationReason string    `json:"cancellationReason"`
	Note               string    `json:"note"`
	OrderDetail        JSONB     `gorm:"type:jsonb;" json:"orderDetail"`
	Deadline           time.Time `json:"deadline"`
	UserID             uuid.UUID `json:"userID" gorm:"type:uuid;"`
	User               User      `json:"user" gorm:"foreignKey:UserID;references:ID"`
	StatusID           uuid.UUID `json:"statusID" gorm:"type:uuid;"`
	Status             Status    `json:"status" gorm:"foreignKey:StatusID;references:ID"`
}
