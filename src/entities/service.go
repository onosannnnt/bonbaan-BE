package Entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type Service struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:(uuid_generate_v4())"`
	Description string    `json:"description"`
	Name        string    `json:"name"`
	Rate        int       `json:"rate"`
	Categories []Category `json:"categories" gorm:"many2many:services_categories;"`
}
