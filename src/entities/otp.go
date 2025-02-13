package Entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Otp struct {
	gorm.Model
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Email   string    `json:"email"`
	Otp     string    `json:"otp"`
	Expired time.Time `json:"expired"`
}
