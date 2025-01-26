package orderUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
)

type OrderUsecase interface {
	Insert(order *Entities.Order) error
	GetAll() ([]*Entities.Order, error)
	GetOne(id *string) (*Entities.Order, error)
	Update(id *string, order *Entities.Order) error
	Delete(id *string) error
}

type OrderService struct {
	orderRepo OrderRepository
}

func NewOrderService(repo OrderRepository) OrderUsecase {
	return &OrderService{
		orderRepo: repo,
	}
}

func (s *OrderService) Insert(order *Entities.Order) error {
	status, err := s.orderRepo.GetDefaultStatus()
	if err != nil {
		return err
	}
	order.Status = *status
	return s.orderRepo.Insert(order)
}

func (s *OrderService) GetAll() ([]*Entities.Order, error) {

	return s.orderRepo.GetAll()
}

func (s *OrderService) GetOne(id *string) (*Entities.Order, error) {
	return s.orderRepo.GetByID(id)
}

func (s *OrderService) Update(id *string, order *Entities.Order) error {
	return s.orderRepo.Update(id, order)
}

func (s *OrderService) Delete(id *string) error {
	return s.orderRepo.Delete(id)
}
