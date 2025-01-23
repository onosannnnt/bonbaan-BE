package Entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JSON json.RawMessage

type Order struct {
	gorm.Model
	ID                 uuid.UUID              `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CancellationReason string                 `json:"cancellationReason"`
	Note               string                 `json:"note"`
	OrderDetail        map[string]interface{} `gorm:"serializer:json;" json:"orderDetail"`
	Dateline           time.Time              `json:"dateline"`
	UserID             uuid.UUID              `json:"userID"`
	User               User                   `gorm:"foreignKey:UserID;references:ID"`
	StatusID           uuid.UUID              `json:"statusID"`
	Status             Status                 `gorm:"foreignKey:StatusID;references:ID"`
}
