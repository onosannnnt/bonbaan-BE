package orderUsecase

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"github.com/onosannnnt/bonbaan-BE/src/config"
	"github.com/onosannnnt/bonbaan-BE/src/constance"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
	serviceUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/service"
	statusUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/status"
)

type OrderUsecase interface {
	Insert(order *Entities.Order) (*Entities.Order, error)
	GetAll(config *model.Pagination) ([]*Entities.Order, *model.Pagination, error)
	GetByID(id *string) (*Entities.Order, error)
	GetByStatus(statusID *uuid.UUID, config *model.Pagination) ([]*Entities.Order, *model.Pagination, error)
	Update(id *string, order *Entities.Order) error
	Delete(id *string) error
	Hook(id *string) error
	CancleOrder(id *string, cancleReason *string) error
	AcceptOrder(id *string) error
	SubmitOrder(order *model.ConfirmOrderRequest) error
	CompleteOrder(id *string) error
}

type OrderService struct {
	orderRepo   OrderRepository
	serviceRepo serviceUsecase.ServiceRepository
	statusRepo  statusUsecase.StatusRepository
}

func NewOrderService(orderRepo OrderRepository, serviceRepo serviceUsecase.ServiceRepository, statusRepo statusUsecase.StatusRepository) OrderUsecase {
	return &OrderService{
		orderRepo:   orderRepo,
		serviceRepo: serviceRepo,
		statusRepo:  statusRepo,
	}
}

func (s *OrderService) Insert(order *Entities.Order) (*Entities.Order, error) {
	status, err := s.orderRepo.GetDefaultStatus()
	if err != nil {
		return nil, err
	}
	client, err := omise.NewClient(config.OmisePublicKey, config.OmiseSecretKey)
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

func (s *OrderService) GetAll(config *model.Pagination) ([]*Entities.Order, *model.Pagination, error) {
	if config.PageSize <= 0 {
		config.PageSize = 10
	}
	if config.CurrentPage <= 0 {
		config.CurrentPage = 1
	}
	orders, totalRecords, err := s.orderRepo.GetAll(config)
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

func (s *OrderService) GetByID(id *string) (*Entities.Order, error) {
	return s.orderRepo.GetByID(id)
}

func (s *OrderService) GetByStatus(statusID *uuid.UUID, config *model.Pagination) ([]*Entities.Order, *model.Pagination, error) {
	if config.PageSize <= 0 {
		config.PageSize = 10
	}
	if config.CurrentPage <= 0 {
		config.CurrentPage = 1
	}
	orders, totalRecords, err := s.orderRepo.GetByStatusID(statusID, config)
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

func (s *OrderService) Update(id *string, order *Entities.Order) error {
	return s.orderRepo.Update(id, order)
}

func (s *OrderService) Delete(id *string) error {
	return s.orderRepo.Delete(id)
}

func (s *OrderService) CancleOrder(id *string, cancleReason *string) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return err
	}
	status, err := s.statusRepo.GetByName(&constance.Status_Cancelled)
	if err != nil {
		return err
	}
	if order.Status.Name == constance.Status_Cancelled {
		return nil
	}
	order.Status = *status
	order.CancellationReason = *cancleReason
	return s.orderRepo.Update(id, order)
}

func (s *OrderService) AcceptOrder(id *string) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return err
	}
	status, err := s.statusRepo.GetByName(&constance.Status_Unpaid)
	if err != nil {
		return err
	}
	order.Status = *status
	client, err := omise.NewClient(config.OmisePublicKey, config.OmiseSecretKey)
	if err != nil {
		return err
	}
	source := &omise.Source{}
	price, err := strconv.ParseFloat(order.OrderDetail["price"].(string), 64)
	if err != nil {
		return err
	}
	err = client.Do(source, &operations.CreateSource{

		Amount:   int64(price * 100),
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
	transaction.Price = price
	transaction.ChargeID = charge.ID
	transaction.Charge = *charge
	order.Transaction = transaction
	order.Status = *status
	err = s.orderRepo.Update(id, order)
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderService) SubmitOrder(order *model.ConfirmOrderRequest) error {
	orderEntity, err := s.orderRepo.GetByID(&order.OrderID)
	if err != nil {
		return err
	}
	status, err := s.statusRepo.GetByName(&constance.Status_Confirm)
	if err != nil {
		return err
	}
	orderEntity.Status = *status
	orderEntity.Attachments = order.Attachments
	return s.orderRepo.Update(&order.OrderID, orderEntity)
}

func (s *OrderService) CompleteOrder(id *string) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return err
	}
	status, err := s.statusRepo.GetByName(&constance.Status_Completed)
	if err != nil {
		return err
	}
	order.Status = *status
	return s.orderRepo.Update(id, order)
}
