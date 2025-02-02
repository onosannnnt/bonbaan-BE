package Entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service_Category struct {
    gorm.Model
    ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:(uuid_generate_v4())"`
    CategoryID  uuid.UUID `gorm:"type:uuid;column:category_id"`
    ServiceID   uuid.UUID `gorm:"type:uuid;column:service_id"`
}