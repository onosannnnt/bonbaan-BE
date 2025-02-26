package Entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct{
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:(uuid_generate_v4())"`
	ServiceID          uuid.UUID              `json:"service_id"`
	
}