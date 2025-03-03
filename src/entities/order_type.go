package Entities

import (
	"github.com/google/uuid"
	"github.com/onosannnnt/bonbaan-BE/src/constance"
	"gorm.io/gorm"
)

type OrderType struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;primaryKey;default:(uuid_generate_v4())"`
	Name string    `json:"name"`
}

func InitializeOrderTypeData(db *gorm.DB) error {
	orderTypes := []OrderType{
		{Name: constance.Types_Vow},
		{Name: constance.Types_Fulfill},
	}
	for _, orderType := range orderTypes {
		if err := db.FirstOrCreate(&orderType, OrderType{Name: orderType.Name}).Error; err != nil {
			return err
		}
	}
	return nil
}
