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
func (d *ReviewDriver) GetAll()([]*Entities.Review,error){
	var reviews []*Entities.Review
	if err := d.db.Find(&reviews).Error;err != nil{
		return nil,err
	}
	return reviews, nil
}

func (d *ReviewDriver) GetByID(id *string)(*Entities.Review,error){
	var review Entities.Review
	if err := d.db.Where("id = ?", id).First(&review).Error; err != nil {
		return nil, err
	}
	return &review, nil

}
func (d *ReviewDriver) Update(id *string,review *Entities.Review) error{
	if err := d.db.Model(review).Updates(review).Error; err != nil {
		return err
	}
	return nil
}

func (d *ReviewDriver) Delete(id *string) error{
	if err := d.db.Where("id = ?", *id).Delete(&Entities.Review{}).Error; err != nil {
		return err
	}
	return nil
}