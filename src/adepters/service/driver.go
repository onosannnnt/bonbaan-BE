package serviceAdapter

import (
	"fmt"

	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
	serviceUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/service"
	"gorm.io/gorm"
)

type ServiceDriver struct {
	db *gorm.DB
}

func NewServiceDriver(db *gorm.DB) serviceUsecase.ServiceRepository {
	return &ServiceDriver{
		db: db,
	}
}

// Implement the Insert method to satisfy the ServiceRepository interface
func (d *ServiceDriver) Insert(service *Entities.Service) error {
	if err := d.db.Create(service).Error; err != nil {
		return err
	}
	return nil
}

func (d *ServiceDriver) GetAll(config *model.Pagination) (*[]Entities.Service, int64, error) {
	var services []Entities.Service
	var totalRecords int64

	db := d.db.Model(&Entities.Service{}).
		Joins("LEFT JOIN review_utils ON review_utils.service_id = services.id")

	if config.Search != "" {
		search := fmt.Sprintf("%%%s%%", config.Search)
		db = db.Where("services.name ILIKE ? OR services.description ILIKE ?", search, search)
	}

	if err := db.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Ordering logic
	if config.OrderBy != "" && (config.OrderDirection == "ASC" || config.OrderDirection == "DESC") {
		if config.OrderBy == "rate" {
			db = db.Order(fmt.Sprintf("review_utils.total_rete / NULLIF(review_utils.total_reviewer, 0) %s", config.OrderDirection))
		} else {
			db = db.Order(fmt.Sprintf("services.%s %s", config.OrderBy, config.OrderDirection))
		}
	} else {
		db = db.Order("services.updated_at DESC")
	}

	if err := db.
		Preload("Categories", func(tx *gorm.DB) *gorm.DB {
			return tx.Omit("created_at", "updated_at", "deleted_at")
		}).
		Preload("Packages", func(tx *gorm.DB) *gorm.DB {
			return tx.Omit("created_at", "updated_at", "deleted_at")
		}).
		Preload("Attachments", func(tx *gorm.DB) *gorm.DB {
			return tx.Omit("created_at", "updated_at", "deleted_at")
		}).
		Limit(config.PageSize).
		Offset((config.CurrentPage - 1) * config.PageSize).
		Find(&services).Error; err != nil {
		return nil, 0, err
	}

	for i := range services {
		var reviewUtils Entities.Review_utils
		if err := d.db.Where("service_id = ?", services[i].ID).First(&reviewUtils).Error; err == nil {
			if reviewUtils.TotalReviewer > 0 {
				services[i].Rate = float64(reviewUtils.TotalRete) / float64(reviewUtils.TotalReviewer)
			}
		}
	}

	return &services, totalRecords, nil
}

func (d *ServiceDriver) GetByID(id *string) (*Entities.Service, error) {
	var service Entities.Service
	if err := d.db.Preload("Review_utils").Preload("Categories").Preload("Packages").Preload("Attachments").Where("id = ?", id).First(&service).Error; err != nil {
		return nil, err
	}

	var reviewUtils Entities.Review_utils
	if err := d.db.Where("service_id = ?", service.ID).First(&reviewUtils).Error; err == nil {
		if reviewUtils.TotalReviewer > 0 {
			service.Rate = float64(reviewUtils.TotalRete) / float64(reviewUtils.TotalReviewer)
		}
	}

	return &service, nil
}

func (d *ServiceDriver) GetPackagebyServiceID(serviceID *string) (*[]Entities.Package, error) {
	var packages []Entities.Package
	if err := d.db.Where("service_id = ?", serviceID).Find(&packages).Error; err != nil {
		return nil, err
	}
	return &packages, nil
}

func (d *ServiceDriver) Update(service *Entities.Service) error {
	if err := d.db.Model(&Entities.Service{}).Where("id = ?", service.ID).Updates(service).Error; err != nil {
		return err
	}
	return nil
}

func (d *ServiceDriver) Delete(id *string) error {
	if err := d.db.Where("id = ?", id).Delete(&Entities.Service{}).Error; err != nil {
		return err
	}
	return nil
}
