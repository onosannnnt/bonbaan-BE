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

// Add this method to the ServiceDriver struct
func (d *ServiceDriver) InitializeFullTextSearchIndex() error {
    // First check if the Thai text search configuration exists
    var thaiConfigExists bool
    err := d.db.Raw("SELECT EXISTS (SELECT 1 FROM pg_ts_config WHERE cfgname = 'thai')").Scan(&thaiConfigExists).Error
    if err != nil {
        return err
    }
    
    // Create Thai text search configuration if it doesn't exist
    if !thaiConfigExists {
        // Create Thai configuration based on simple
        err = d.db.Exec("CREATE TEXT SEARCH CONFIGURATION thai (COPY = simple)").Error
        if err != nil {
            return err
        }
        
        // Alter the mapping to use simple dictionary for word type
        err = d.db.Exec("ALTER TEXT SEARCH CONFIGURATION thai ALTER MAPPING FOR word WITH simple").Error
        if err != nil {
            return err
        }
    }
    
    // Check if the index already exists
    var indexExists bool
    err = d.db.Raw("SELECT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_services_fts')").Scan(&indexExists).Error
    if err != nil {
        return err
    }

    // Create the index using the Thai configuration if it doesn't exist
    if !indexExists {
        return d.db.Exec("CREATE INDEX idx_services_fts ON services USING gin(to_tsvector('thai', name || ' ' || description || ' ' || address))").Error
    }
	//Ensure pg_trgm is enable
	err = d.db.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm").Error
	if err != nil {
		return err
	}

    
    return nil
}

// Update the NewServiceDriver function to initialize the index
func NewServiceDriver(db *gorm.DB) serviceUsecase.ServiceRepository {
    driver := &ServiceDriver{
        db: db,
    }
    
    // Initialize the full-text search index (ignore error for simplicity)
    _ = driver.InitializeFullTextSearchIndex()
    
    return driver
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

	searchQuery := config.Search // สมมุติว่ามีค่าจาก user เช่น "ไอ้ไข่"

	if searchQuery != "" {
		db = db.Select(`
			services.*,
			ts_rank(
				setweight(to_tsvector('thai', coalesce(services.name, '')), 'A') || 
				setweight(to_tsvector('thai', coalesce(services.address, '')), 'C') ||
				setweight(to_tsvector('thai', coalesce(services.description, '')), 'B'),
				plainto_tsquery('thai', ?)
			) AS rank,
			similarity(services.name || ' ' || services.description || ' ' || services.address, ?) AS sim
		`, searchQuery, searchQuery).
			Where(`
				(
					(
						setweight(to_tsvector('thai', coalesce(services.name, '')), 'A') || 
						setweight(to_tsvector('thai', coalesce(services.address, '')), 'C') ||
						setweight(to_tsvector('thai', coalesce(services.description, '')), 'B')
					) @@ plainto_tsquery('thai', ?)
					OR similarity(services.name || ' ' || services.description || ' ' || services.address, ?) > 0.0
				)
			`, searchQuery, searchQuery).
			Order("rank DESC, sim DESC")
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
	if err := d.db.Preload("Review_utils").Preload("Categories").Preload("Packages").Preload("Packages.OrderType").Preload("Attachments").Where("id = ?", id).First(&service).Error; err != nil {
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
	tx := d.db.Begin()
	if err := tx.Model(&Entities.Service{}).Where("id = ?", service.ID).Updates(service).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update Categories association only if provided.
	// fmt.Println(service.Categories)
	if service.Categories != nil {
		if err := tx.Model(service).Association("Categories").Replace(service.Categories); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update Packages association only if provided.
	// fmt.Println(service.Packages)
	if service.Packages != nil {
		if err := tx.Model(service).Association("Packages").Replace(service.Packages); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update Attachments association only if provided.
	if service.Attachments != nil {
		if err := tx.Model(service).Association("Attachments").Replace(service.Attachments); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (d *ServiceDriver) Delete(id *string) error {
	if err := d.db.Where("id = ?", id).Delete(&Entities.Service{}).Error; err != nil {
		return err
	}
	return nil
}
