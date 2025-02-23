package Entities

import (
	"github.com/google/uuid"
	"github.com/omise/omise-go"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	ID       uuid.UUID    `gorm:"primaryKey;default:(uuid_generate_v4())"`
	Price    float64      `json:"price"`
	ChargeID string       `json:"chargeID"`
	Charge   omise.Charge `gorm:"serializer:json;" json:"charge"`
}
