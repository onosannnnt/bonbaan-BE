package Entities

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Package struct {
	gorm.Model
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:(uuid_generate_v4())"`
	Name        string         `json:"name"`
	Item        pq.StringArray `gorm:"type:text[]" json:"item"`
	Price       int            `json:"price"`
	Description string         `json:"description"`
	ServiceID   uuid.UUID      `json:"service_id"`
	OrderTypeID uuid.UUID      `json:"order_type_id"`
	OrderType   OrderType      `gorm:"foreignKey:OrderTypeID ;references:ID"`
	// Service     Service `gorm:"foreignKey:ServiceID ;references:ID"`

}
