package serviceUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
)

type ServiceUsecase interface {
	CreateService(service *Entities.Service) error
	GetAll(config *model.Pagination) (*[]Entities.Service, *model.Pagination, error)
	GetByID(id *string) (*Entities.Service, error)
	GetPackageByServiceID(serviceID *string) (*[]Entities.Package, error) // New method to get presets by service ID
	UpdateService(service *Entities.Service) error
	DeleteService(id *string) error
}

type ServiceAsService struct {
	ServiceRepo ServiceRepository
}

func NewServiceUsecase(repo ServiceRepository) ServiceUsecase {
	return &ServiceAsService{ServiceRepo: repo}
}

func (sc *ServiceAsService) CreateService(service *Entities.Service) error {
	return sc.ServiceRepo.Insert(service)
}

// Implement the GetAll method to satisfy the ServiceUsecase interface
func (sc *ServiceAsService) GetAll(config *model.Pagination) (*[]Entities.Service, *model.Pagination, error) {
	if config.PageSize <= 0 {
		config.PageSize = 10
	}
	if config.CurrentPage <= 0 {
		config.CurrentPage = 1
	}
	orders, totalRecords, err := sc.ServiceRepo.GetAll(config)
	if err != nil {
		return nil, nil, err
	}
	var totalPages int64
	if totalRecords%int64(config.PageSize) == 0 {
		totalPages = totalRecords / int64(config.PageSize)
	} else {
		totalPages = totalRecords/int64(config.PageSize) + 1
	}
	pagination := &model.Pagination{
		CurrentPage:  config.CurrentPage,
		PageSize:     config.PageSize,
		TotalPages:   int(totalPages),
		TotalRecords: int(totalRecords),
	}

	return orders, pagination, nil
}

func (sc *ServiceAsService) GetByID(id *string) (*Entities.Service, error) {
	// Implementation of GetByID method
	return sc.ServiceRepo.GetByID(id)
}
func (sc *ServiceAsService) GetPackageByServiceID(serviceID *string) (*[]Entities.Package, error) {
	// Implementation of GetPackagessbyServiceID method
	return sc.ServiceRepo.GetPackagebyServiceID(serviceID)
}

func (sc *ServiceAsService) UpdateService(service *Entities.Service) error {
	// Implementation of UpdateService method
	return sc.ServiceRepo.Update(service)
}

func (sc *ServiceAsService) DeleteService(id *string) error {
	// Implementation of DeleteService method
	return sc.ServiceRepo.Delete(id)
}
