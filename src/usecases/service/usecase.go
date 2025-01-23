package serviceUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
)

type ServiceUsecase interface {
    CreateService(service *Entities.Service) error
    GetAll() (*[]Entities.Service, error)
	GetByID(id *string) (*Entities.Service, error)
	UpdateService(service *Entities.Service) error


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
func (sc *ServiceAsService) GetAll() (*[]Entities.Service, error) {
    // Implementation of GetAll method
    return sc.ServiceRepo.GetAll()
}

func (sc *ServiceAsService) GetByID(id *string) (*Entities.Service, error) {
    // Implementation of GetByID method
    return sc.ServiceRepo.GetByID(id)
}

func (sc *ServiceAsService) UpdateService(service *Entities.Service) error {
	// Implementation of UpdateService method
	return sc.ServiceRepo.Update(service)
}
