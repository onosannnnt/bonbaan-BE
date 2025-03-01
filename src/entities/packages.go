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
	PackageTypeID uuid.UUID `json:"package_type_id"`
	PackageType PackageType `gorm:"foreignKey:PackageTypeID ;references:ID"` 
	// Service     Service `gorm:"foreignKey:ServiceID ;references:ID"`

}
