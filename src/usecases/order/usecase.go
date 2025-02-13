package orderUsecase

import (
	"strconv"

	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"github.com/onosannnnt/bonbaan-BE/src/Config"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	serviceUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/service"
)

type OrderUsecase interface {
	Insert(order *Entities.Order) (*Entities.Order, error)
	GetAll() ([]*Entities.Order, error)
	GetByID(id *string) (*Entities.Order, error)
	Update(id *string, order *Entities.Order) error
	Delete(id *string) error
	Hook(id *string) error
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

func (s *OrderService) Insert(order *Entities.Order) (*Entities.Order, error) {
	status, err := s.orderRepo.GetDefaultStatus()
	if err != nil {
		return nil, err
	}
	client, err := omise.NewClient(Config.OmisePublicKey, Config.OmiseSecretKey)
	if err != nil {
		return nil, err
	}
	source := &omise.Source{}
	price, err := strconv.ParseFloat(order.OrderDetail["price"].(string), 64)
	if err != nil {
		return nil, err
	}
	err = client.Do(source, &operations.CreateSource{

		Amount:   int64(price * 100),
		Currency: "thb",
		Type:     "promptpay",
	})
	if err != nil {
		return nil, err
	}
	charge := &omise.Charge{}
	err = client.Do(charge, &operations.CreateCharge{
		Amount:   source.Amount,
		Currency: source.Currency,
		Source:   source.ID,
	})
	if err != nil {
		return nil, err
	}
	var transaction Entities.Transaction
	transaction.Price = price
	transaction.ChargeID = charge.ID
	transaction.Charge = *charge
	order.Transaction = transaction
	order.Status = *status
	err = s.orderRepo.Insert(order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) Hook(ChargeID *string) error {
	err := s.orderRepo.GetAndUpdateByChargeID(*ChargeID)
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderService) GetAll() ([]*Entities.Order, error) {

	return s.orderRepo.GetAll()
}

func (s *OrderService) GetByID(id *string) (*Entities.Order, error) {
	return s.orderRepo.GetByID(id)
}

func (s *OrderService) Update(id *string, order *Entities.Order) error {
	return s.orderRepo.Update(id, order)
}

func (s *OrderService) Delete(id *string) error {
	return s.orderRepo.Delete(id)
}
