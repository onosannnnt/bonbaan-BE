package orderUsecase

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"github.com/onosannnnt/bonbaan-BE/src/config"
	"github.com/onosannnnt/bonbaan-BE/src/constance"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
	NotificationUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/notification"
	serviceUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/service"
	statusUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/status"
	"gorm.io/gorm"
)

type OrderUsecase interface {
	Insert(order *Entities.Order) (*Entities.Order, error)
	GetAll(config *model.Pagination) ([]*Entities.Order, *model.Pagination, error)
	GetByID(id *string) (*Entities.Order, error)
	GetByStatus(statusID *uuid.UUID, config *model.Pagination) ([]*Entities.Order, *model.Pagination, error)
	Update(id *string, order *Entities.Order) error
	Delete(id *string) error
	Hook(id *string) error
	CancelOrder(id *string, cancelReason *string) error
	SubmitOrder(order *model.SubmitOrderRequest) error
	CompleteOrder(id *string) error
	InsertCustomOrder(order *Entities.Order) (*Entities.Order, error)
	AcceptOrder(data *model.ConfirmOrderRequest) error
	ApproveOrder(id *string) error
}

type OrderService struct {
	db               *gorm.DB
	orderRepo        OrderRepository
	serviceRepo      serviceUsecase.ServiceRepository
	statusRepo       statusUsecase.StatusRepository
	notificationRepo NotificationUsecase.NotificationRepository
}

func NewOrderService(orderRepo OrderRepository, serviceRepo serviceUsecase.ServiceRepository, statusRepo statusUsecase.StatusRepository, notificationRepo NotificationUsecase.NotificationRepository, db *gorm.DB) OrderUsecase {
	return &OrderService{
		db:               db,
		orderRepo:        orderRepo,
		serviceRepo:      serviceRepo,
		statusRepo:       statusRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *OrderService) Insert(order *Entities.Order) (*Entities.Order, error) {
	status, err := s.statusRepo.GetByName(&constance.Status_Unpaid)
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
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
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
		tx.Rollback()
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
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
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

func (s *OrderService) CancelOrder(id *string, cancelReason *string) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	status, err := s.statusRepo.GetByName(&constance.Status_Cancelled)
	if err != nil {
		tx.Rollback()
		return err
	}
	if order.Status.Name == constance.Status_Cancelled {
		return errors.New("order is already cancelled")
	}
	order.Status = *status
	err = s.notificationRepo.Insert(&Entities.Notification{
		UserID:  order.UserID,
		Header:  "Order Cancelled",
		Body:    "Your order has been cancelled because " + *cancelReason,
		OrderID: order.ID,
	})
	if err != nil {
		tx.Rollback()
	}
	s.orderRepo.Update(id, order)
	tx.Commit()
	return nil
}

func (s *OrderService) ApproveOrder(id *string) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err

	}
	if order.Status.Name != constance.Status_Confirm {
		return errors.New("order is not confirm")
	}
	status, err := s.statusRepo.GetByName(&constance.Status_Unpaid)
	if err != nil {
		tx.Rollback()
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
		tx.Rollback()
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
		tx.Rollback()
		return err
	}
	err = s.notificationRepo.Insert(&Entities.Notification{
		UserID:  order.UserID,
		Header:  "Order Approved",
		Body:    "Your order has been approved",
		OrderID: order.ID,
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (s *OrderService) SubmitOrder(order *model.SubmitOrderRequest) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	orderEntity, err := s.orderRepo.GetByID(&order.OrderID)
	if err != nil {
		tx.Rollback()
		return err
	}
	if orderEntity.Status.Name != constance.Status_Processing {
		tx.Rollback()
		return errors.New("order is not in processing")
	}
	status, err := s.statusRepo.GetByName(&constance.Status_Approve)
	if err != nil {
		tx.Rollback()
		return err
	}
	orderEntity.Status = *status
	orderEntity.Attachments = order.Attachments
	s.orderRepo.Update(&order.OrderID, orderEntity)
	err = s.notificationRepo.Insert(&Entities.Notification{
		UserID:  orderEntity.UserID,
		Header:  "Order Submitted",
		Body:    fmt.Sprintf("Your %s order has been submitted. Please check your order to approve", orderEntity.ID),
		OrderID: orderEntity.ID,
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (s *OrderService) CompleteOrder(id *string) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return err
	}
	status, err := s.statusRepo.GetByName(&constance.Status_Review)
	if err != nil {
		return err
	}
	order.Status = *status
	s.orderRepo.Update(id, order)
	err = s.notificationRepo.Insert(&Entities.Notification{
		UserID:  order.UserID,
		Header:  "Order Completed",
		Body:    "Your order has been completed",
		OrderID: order.ID,
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (s *OrderService) InsertCustomOrder(order *Entities.Order) (*Entities.Order, error) {
	status, err := s.statusRepo.GetByName(&constance.Status_Pending)
	if err != nil {
		return nil, err
	}
	order.Status = *status
	err = s.orderRepo.Insert(order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) AcceptOrder(data *model.ConfirmOrderRequest) error {
	order, err := s.orderRepo.GetByID(&data.OrderID)
	if err != nil {
		return err
	}
	if order.Status.Name != constance.Status_Pending {
		return errors.New("order is not pending")
	}
	status, err := s.statusRepo.GetByName(&constance.Status_Confirm)
	if err != nil {
		return err
	}
	order.Status = *status
	order.OrderDetail["price"] = data.Price
	return s.orderRepo.Update(&data.OrderID, order)
}
