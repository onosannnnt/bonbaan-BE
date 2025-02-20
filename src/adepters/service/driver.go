package serviceAdapter

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
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

func (d *ServiceDriver) GetAll() (*[]Entities.Service, error) {
	var services []Entities.Service
	if err := d.db.Find(&services).Error; err != nil {
		return nil, err
	}
	return &services, nil
}
func (d *ServiceDriver) GetByID(id *string) (*Entities.Service, error) {
	var service Entities.Service
	if err := d.db.Where("id = ?", id).First(&service).Error; err != nil {
		return nil, err
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
