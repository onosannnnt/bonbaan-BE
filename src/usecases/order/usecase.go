package orderUsecase

import (
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"github.com/onosannnnt/bonbaan-BE/src/Config"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	serviceUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/service"
)

type OrderUsecase interface {
	Insert(order *Entities.Order) error
	GetAll() ([]*Entities.Order, error)
	GetOne(id *string) (*Entities.Order, error)
	Update(id *string, order *Entities.Order) error
	Delete(id *string) error
}

type OrderService struct {
	orderRepo   OrderRepository
	serviceRepo serviceUsecase.ServiceRepository
}

func NewOrderService(orderRepo OrderRepository, serviceRepo serviceUsecase.ServiceRepository) OrderUsecase {
	return &OrderService{
		orderRepo:   orderRepo,
		serviceRepo: serviceRepo,
	}
}

func (s *OrderService) Insert(order *Entities.Order) error {
	status, err := s.orderRepo.GetDefaultStatus()
	if err != nil {
		return err
	}
	client, err := omise.NewClient(Config.OmisePublicKey, Config.OmiseSecretKey)
	if err != nil {
		return err
	}
	source := &omise.Source{}
	err = client.Do(source, &operations.CreateSource{
		Amount:   int64(20 * 100),
		Currency: "thb",
		Type:     "promptpay",
	})
	if err != nil {
		return err
	}
	charge := &omise.Charge{}
	err = client.Do(charge, &operations.CreateCharge{
		Amount:   source.Amount,
		Currency: source.Currency,
		Source:   source.ID,
	})
	if err != nil {
		return err
	}
	var transaction Entities.Transaction
	transaction.Price = 20
	transaction.Charge = *charge
	order.Transaction = transaction
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
