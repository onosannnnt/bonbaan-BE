package Entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	ID           uuid.UUID    `gorm:"type:uuid;primaryKey;default:(uuid_generate_v4())"`
	Description  string       `json:"description"`
	Name         string       `json:"name"`
	Rate         float64      `json:"rate" gorm:"default:0.0"`
	Review_utils Review_utils `gorm:"foreignKey:ServiceID"`
	Address      string       `json:"address"`
	Categories   []Category   `json:"categories" gorm:"many2many:services_categories;"`
	Packages     []Package    `json:"packages" gorm:"foreignKey:ServiceID"`
	Attachments  []Attachment `json:"attachments" gorm:"foreignKey:ServiceID"`
}
