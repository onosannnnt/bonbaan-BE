package Entities

import (
	"github.com/google/uuid"
	"github.com/onosannnnt/bonbaan-BE/src/constance"
	"gorm.io/gorm"
)

type Status struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:(uuid_generate_v4())"`
	Name string    `json:"name"`
}

func InitializeStatusData(db *gorm.DB) error {
	statuses := []Status{
		{Name: constance.Status_Pending},
		{Name: constance.Status_Unpaid},
		{Name: constance.Status_Processing},
		{Name: constance.Status_Confirm},
		{Name: constance.Status_Review},
		{Name: constance.Status_Completed},
		{Name: constance.Status_Refund},
		{Name: constance.Status_Cancelled},
	}
	for _, status := range statuses {
		if err := db.FirstOrCreate(&status, Status{Name: status.Name}).Error; err != nil {
			return err
		}
	}
	return nil
}
