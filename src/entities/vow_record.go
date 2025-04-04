package Entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VowRecord struct {
	gorm.Model
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;default:(uuid_generate_v4())"`
	Vow              string    `json:"vow"`
	Deadline         time.Time `json:"deadline"`
	Note             string    `json:"note"`
	UserID           uuid.UUID `json:"user_id" gorm:"foreignKey:UserID;default:null"`
	User             User      `json:"user" gorm:"foreignKey:UserID"`
	ServiceID        uuid.UUID `json:"service_id" gorm:"foreignKey:ServiceID;default:null"`
	Service          Service   `json:"service" gorm:"foreignKey:ServiceID"`
	VowOrderID       uuid.UUID `json:"vow_order_id" gorm:"foreignKey:VowOrderID;default:null"`
	VowOrder         Order     `json:"vow_order" gorm:"foreignKey:VowOrderID"`
	FulfilledOrderID uuid.UUID `json:"fulfilled" gorm:"default:null"`
	FulfillOrder     Order     `json:"fulfill_order" gorm:"foreignKey:FulfilledOrderID"`
}
