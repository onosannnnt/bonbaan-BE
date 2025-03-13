package Entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/onosannnnt/bonbaan-BE/src/utils"
	"gorm.io/gorm"
)

type JSON json.RawMessage

type Order struct {
	gorm.Model
	ID                 string                 `gorm:"primaryKey)"`
	CancellationReason string                 `json:"cancellationReason"`
	Note               string                 `json:"note"`
	OrderDetail        map[string]interface{} `gorm:"serializer:json;" json:"orderDetail"`
	Deadline           time.Time              `json:"deadline"`
	UserID             uuid.UUID              `json:"userID"`
	User               User                   `gorm:"foreignKey:UserID;references:ID"`
	StatusID           uuid.UUID              `json:"statusID"`
	Status             Status                 `gorm:"foreignKey:StatusID;references:ID"`
	ServiceID          uuid.UUID              `json:"serviceID"`
	Service            Service                `gorm:"foreignKey:ServiceID;references:ID"`
	TransactionID      uuid.UUID              `json:"transactionID" gorm:"foreignKey:TransactionID;references:ID;default:null"`
	Transaction        Transaction            `gorm:"foreignKey:TransactionID;references:ID;default:null"`
	Attachments        []Attachment           `json:"attachments" gorm:"foreignKey:OrderID"`
	OrderType          OrderType              `json:"OrderType"`
	OrderTypeID        uuid.UUID              `gorm:"foreignKey:OrderTypeID;references:ID" json:"OrderTypeID"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID, err = utils.GenerateRandomID()
	var existingOrder Order
	result := tx.Find(&existingOrder, "id = ?", o.ID)
	if result.RowsAffected > 0 {
		o.BeforeCreate(tx)
	}
	if err != nil {
		return err
	}
	return
}
