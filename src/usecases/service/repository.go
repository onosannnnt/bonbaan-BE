package serviceUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
)

type ServiceRepository interface {
	Insert(ser *Entities.Service) error
	GetAll() (*[]Entities.Service, error)
	GetByID(id *string) (*Entities.Service, error)
	Update(service *Entities.Service) error
	// Delete(id *string) error
	
}
