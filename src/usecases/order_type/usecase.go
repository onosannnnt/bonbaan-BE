package orderTypeUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
)

type OrderTypeUsecase interface {
	Insert(orderType *Entities.OrderType) error
	GetAll() (*[]Entities.OrderType, error)
	GetByID(id *string) (*Entities.OrderType, error)
	Update(id *string, orderType *Entities.OrderType) error
	Delete(id *string) error
}

type OrderTypeService struct {
	orderTypeRepo OrderTypeRepository
}

func NewOrderTypeService(repo OrderTypeRepository) OrderTypeUsecase {
	return &OrderTypeService{
		orderTypeRepo: repo,
	}
}

func (s *OrderTypeService) Insert(orderType *Entities.OrderType) error {
	return s.orderTypeRepo.Insert(orderType)
}

func (s *OrderTypeService) GetAll() (*[]Entities.OrderType, error) {
	return s.orderTypeRepo.GetAll()
}

func (s *OrderTypeService) GetByID(id *string) (*Entities.OrderType, error) {
	return s.orderTypeRepo.GetByID(id)
}

func (s *OrderTypeService) Update(id *string, orderType *Entities.OrderType) error {
	return s.orderTypeRepo.Update(id, orderType)
}

func (s *OrderTypeService) Delete(id *string) error {
	return s.orderTypeRepo.Delete(id)
}
