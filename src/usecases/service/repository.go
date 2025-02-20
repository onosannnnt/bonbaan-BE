package serviceUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
)

type ServiceRepository interface {
	Insert(ser *Entities.Service) error
	GetAll(config *model.Pagination) (*[]Entities.Service, int64, error)
	GetByID(id *string) (*Entities.Service, error)
	Update(service *Entities.Service) error
	Delete(id *string) error
}
