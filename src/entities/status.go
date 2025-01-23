package Entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Status struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name string    `json:"name"`
}
