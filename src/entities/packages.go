package Entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Package struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:(uuid_generate_v4())"`
	Name        string  `json:"name"`
	Item        string  `json:"item"`
	Price       int `json:"price"`
	Description string  `json:"description"`
	ServiceID   uuid.UUID  `json:"service_id"`
	Service     Service `gorm:"foreignKey:ServiceID ;references:ID"`

}
