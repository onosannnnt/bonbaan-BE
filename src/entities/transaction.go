package Entities

import (
	"github.com/google/uuid"
	"github.com/omise/omise-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Transaction struct {
	gorm.Model
	ID       uuid.UUID    `gorm:"primaryKey;default:(uuid_generate_v4())"`
	Price    float64      `json:"price"`
	ChargeID string       `json:"chargeID"`
	Charge   omise.Charge `gorm:"serializer:json;" json:"charge"`
}

func (t *Transaction) AfterDelete(tx *gorm.DB) (err error) {
	if err := tx.Clauses(clause.Returning{}).Where("ID = ?", t.ID).Delete(&Order{}).Error; err != nil {
		return err
	}
	return nil
}
