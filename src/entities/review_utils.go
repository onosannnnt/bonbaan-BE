package Entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review_utils struct{
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:(uuid_generate_v4())"`
	ServiceID          uuid.UUID              `json:"service_id"`
	TotalRete	int  `json:"total_rete" gorm:"default:0"`
	TotalReviewer	int  `json:"total_reviewer" gorm:"default:0"`
}
