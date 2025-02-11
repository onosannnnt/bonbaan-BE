package reviewAdepter

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	reviewUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/review"
	"gorm.io/gorm"
)

type ReviewDriver struct {
	db *gorm.DB
}

func NewReviewDriver(db *gorm.DB) reviewUsecase.ReviewRepository {
	return &ReviewDriver{
		db: db,
	}
}

func (d *ReviewDriver) Insert(review *Entities.Review) error {
	if err := d.db.Create(review).Error; err != nil{
		return err
	}
	return nil
}