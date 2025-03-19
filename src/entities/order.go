package Entities

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type JSON json.RawMessage

type Order struct {
	gorm.Model
	ID                 uuid.UUID      `gorm:"type:uuid;primaryKey;default:(uuid_generate_v4())"`
	UserID             uuid.UUID      `json:"user_id" gorm:"foreignKey:UserID;default:null"`
	User               User           `json:"user" gorm:"foreignKey:UserID"`
	Price              float64        `json:"price"`
	Items              pq.StringArray `gorm:"type:text[]" json:"items"`
	PackageID          uuid.UUID      `json:"package_id" gorm:"foreignKey:PackageID;default:null"`
	Package            Package        `json:"package" gorm:"foreignKey:PackageID"`
	TransactionID      uuid.UUID      `json:"transaction_id" gorm:"foreignKey:TransactionID;default:null"`
	Transaction        Transaction    `json:"transaction" gorm:"foreignKey:TransactionID"`
	CancellationReason string         `json:"cancellation_reason" gorm:"default:null"`
	StatusID           uuid.UUID      `json:"status_id" gorm:"foreignKey:StatusID;default:null"`
	Status             Status         `json:"status" gorm:"foreignKey:StatusID"`
	Attachments        []Attachment   `json:"attachment" gorm:"foreignKey:OrderID"`
}
