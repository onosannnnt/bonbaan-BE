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
	if err := d.db.Create(review).Error; err != nil {
		return err
	}

	var reviewUtils Entities.Review_utils
	if err := d.db.Where("service_id = ?", review.ServiceID).First(&reviewUtils).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			reviewUtils = Entities.Review_utils{
				ServiceID:     review.ServiceID,
				TotalRete:     review.Rating,
				TotalReviewer: 1,
			}
			if err := d.db.Create(&reviewUtils).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		reviewUtils.TotalRete += review.Rating
		reviewUtils.TotalReviewer += 1
		if err := d.db.Save(&reviewUtils).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d *ReviewDriver) GetAll() ([]*Entities.Review, error) {
	var reviews []*Entities.Review
	if err := d.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Omit("password")
	}).Preload("Service").Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

func (d *ReviewDriver) GetByID(id *string) (*Entities.Review, error) {
	var review Entities.Review
	if err := d.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Omit("password")
	}).Preload("Service").Where("id = ?", id).First(&review).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

func (d *ReviewDriver) Update(id *string, review *Entities.Review) error {
	var oldReview Entities.Review
	if err := d.db.Where("id = ?", id).First(&oldReview).Error; err != nil {
		return err
	}

	if err := d.db.Model(review).Updates(review).Error; err != nil {
		return err
	}

	var reviewUtils Entities.Review_utils
	if err := d.db.Where("service_id = ?", review.ServiceID).First(&reviewUtils).Error; err != nil {
		return err
	}

	reviewUtils.TotalRete = reviewUtils.TotalRete - oldReview.Rating + review.Rating

	if err := d.db.Save(&reviewUtils).Error; err != nil {
		return err
	}

	return nil
}

func (d *ReviewDriver) Delete(id *string) error {
	var review Entities.Review
	if err := d.db.Where("id = ?", *id).First(&review).Error; err != nil {
		return err
	}

	if err := d.db.Where("id = ?", *id).Delete(&Entities.Review{}).Error; err != nil {
		return err
	}

	var reviewUtils Entities.Review_utils
	if err := d.db.Where("service_id = ?", review.ServiceID).First(&reviewUtils).Error; err != nil {
		return err
	}

	reviewUtils.TotalRete -= review.Rating
	reviewUtils.TotalReviewer -= 1

	if err := d.db.Save(&reviewUtils).Error; err != nil {
		return err
	}

	return nil
}
