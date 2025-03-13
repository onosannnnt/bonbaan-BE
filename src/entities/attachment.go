package Entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Attachment struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:(uuid_generate_v4())"`
	URL       string    `json:"url"`
	ServiceID uuid.UUID `json:"service_id" gorm:"foreignKey:ServiceID;default:null"`
	Service   Service   `json:"service" gorm:"foreignKey:ServiceID;default:null"`
	OrderID   string    `json:"order_id" gorm:"foreignKey:ServiceID;default:null"`
	Order     Order     `json:"order" gorm:"foreignKey:OrderID;default:null"`
}
