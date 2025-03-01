package Entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PackageType struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;primaryKey;default:(uuid_generate_v4())"`
	Name string    `json:"name"`
}
