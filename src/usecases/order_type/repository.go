package orderTypeUsecase

import Entities "github.com/onosannnnt/bonbaan-BE/src/entities"

type OrderTypeRepository interface {
	Insert(role *Entities.OrderType) error
	GetAll() (*[]Entities.OrderType, error)
	GetByID(id *string) (*Entities.OrderType, error)
	Update(id *string, role *Entities.OrderType) error
	Delete(id *string) error
}
