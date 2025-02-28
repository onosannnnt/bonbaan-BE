package Entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Attachment struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:(uuid_generate_v4())"`
	URL       string    `json:"url"`
	ServiceID uuid.UUID `json:"service_id"`
	Service   Service   `json:"service" gorm:"foreignKey:ServiceID"`
	OrderID   uuid.UUID `json:"order_id"`
	Order     Order     `json:"order" gorm:"foreignKey:OrderID"`
}
