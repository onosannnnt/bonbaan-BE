package Entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResetPassword struct {
	gorm.Model
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Email         string    `json:"email"`
	ResetPassword string    `json:"reset_password"`
	Expired       time.Time `json:"expired"`
}
