package Entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Recommendation struct {
	gorm.Model
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Current_service_id uuid.UUID `json:"current_service_id" gorm:"type:uuid;not null"`
	Next_service_id    uuid.UUID `json:"next_service_id " gorm:"type:uuid;not null"`
	Total int `json:"total"`
}

type RecommendationUtil struct {
	gorm.Model
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Current_service_id uuid.UUID `json:"current_service_id " gorm:"type:uuid;not null"`
	Total int `json:"total"`
}


