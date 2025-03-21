package orderUsecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"github.com/onosannnnt/bonbaan-BE/src/config"
	"github.com/onosannnnt/bonbaan-BE/src/constance"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
	orderTypeUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/order_type"
	packageUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/package"
	serviceUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/service"
	statusUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/status"
	vowRecordUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/vow_record"
	"gorm.io/gorm"
)

type OrderUsecase interface {
	Insert(order *model.OrderInputRequest) error
	GetAll(config *model.Pagination) ([]*Entities.Order, *model.Pagination, error)
	GetByID(id *string) (*Entities.Order, error)
	GetByStatus(statusID *uuid.UUID, config *model.Pagination) ([]*Entities.Order, *model.Pagination, error)
	Update(id *string, order *Entities.Order) error
	Delete(id *string) error
	Hook(id *string) error
	CancelOrder(id *string, cancelReason *string) error
	SubmitOrder(order *model.SubmitOrderRequest) error
	CompleteOrder(id *string) error
	InsertCustomOrder(order *model.OrderInputRequest) error
	AcceptOrder(data *model.ConfirmOrderRequest) error
	ApproveOrder(id *string) error
}

type OrderService struct {
	db            *gorm.DB
	orderRepo     OrderRepository
	serviceRepo   serviceUsecase.ServiceRepository
	statusRepo    statusUsecase.StatusRepository
	packageRepo   packageUsecase.PackageRepository
	vowRecordRepo vowRecordUsecase.VowRecordRepository
	orderTypeRepo orderTypeUsecase.OrderTypeRepository
}

func NewOrderService(orderRepo OrderRepository, serviceRepo serviceUsecase.ServiceRepository, statusRepo statusUsecase.StatusRepository, db *gorm.DB, packageRepo packageUsecase.PackageRepository, vowRecordRepo vowRecordUsecase.VowRecordRepository, orderTypeRepo orderTypeUsecase.OrderTypeRepository) OrderUsecase {
	return &OrderService{
		db:            db,
		orderRepo:     orderRepo,
		serviceRepo:   serviceRepo,
		statusRepo:    statusRepo,
		packageRepo:   packageRepo,
		vowRecordRepo: vowRecordRepo,
		orderTypeRepo: orderTypeRepo,
	}
}

func (s *OrderService) Insert(order *model.OrderInputRequest) error {
	status, err := s.statusRepo.GetByName(&constance.Status_Unpaid)
	if err != nil {
		return err
	}
	packages, err := s.packageRepo.GetByID(&order.PackageID)
	if err != nil {
		return err
	}
	if packages == nil {
		return errors.New("package not found")
	}
	client, err := omise.NewClient(config.OmisePublicKey, config.OmiseSecretKey)
	if err != nil {
		return err
	}
	source := &omise.Source{}
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err = client.Do(source, &operations.CreateSource{
		Amount:   int64(order.Price * 100),
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
	transaction.Price = order.Price
	transaction.ChargeID = charge.ID
	transaction.Charge = *charge
	var orderEntity Entities.Order
	orderEntity.Price = order.Price
	orderEntity.Package = *packages
	orderEntity.Items = packages.Item
	orderEntity.UserID = uuid.MustParse(order.UserID)
	orderEntity.Status = *status
	orderEntity.ServiceID = packages.ServiceID
	orderEntity.Transaction = transaction
	err = s.orderRepo.Insert(&orderEntity)
	if err != nil {
		tx.Rollback()
		return err
	}
	orderType, err := s.orderTypeRepo.GetByID(&order.OrderTypeID)
	if err != nil {
		return err
	}
	if orderType.Name == constance.Types_Vow {
		parsedDate, err := time.Parse("2006-01-02", order.Deadline)
		if err != nil {
			return err
		}
		vowRecord := Entities.VowRecord{
			Vow:        order.Vow,
			Deadline:   time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.UTC),
			Note:       order.Note,
			UserID:     uuid.MustParse(order.UserID),
			ServiceID:  packages.ServiceID,
			VowOrderID: orderEntity.ID,
		}
		err = s.vowRecordRepo.Insert(&vowRecord)
		if err != nil {
			return errors.New("failed to insert vow record")
		}
	} else if orderType.Name == constance.Types_Fulfill {
		vowRecord := Entities.VowRecord{
			FulfilledOrderID: orderEntity.ID,
		}
		err = s.vowRecordRepo.Update(&order.VowRecordID, &vowRecord)
		if err != nil {
			return errors.New("failed to update vow record")
		}
	}
	tx.Commit()
	return nil
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
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return err
	}
	status, err := s.statusRepo.GetByName(&constance.Status_Cancelled)
	if err != nil {
		return err
	}
	if order.Status.Name == constance.Status_Cancelled {
		return errors.New("order is already cancelled")
	}
	order.Status = *status
	order.CancellationReason = *cancelReason
	return s.orderRepo.Update(id, order)
}

func (s *OrderService) ApproveOrder(id *string) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return err
	}
	if order.Status.Name != constance.Status_Confirm {
		return errors.New("order is not confirm")
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
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err = client.Do(source, &operations.CreateSource{

		Amount:   int64(order.Price * 100),
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
	transaction.Price = order.Price
	transaction.ChargeID = charge.ID
	transaction.Charge = *charge
	order.Transaction = transaction
	order.Status = *status
	err = s.orderRepo.Update(id, order)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *OrderService) SubmitOrder(order *model.SubmitOrderRequest) error {
	orderEntity, err := s.orderRepo.GetByID(&order.OrderID)
	if err != nil {
		return err
	}
	if orderEntity.Status.Name != constance.Status_Processing {
		return errors.New("order is not in processing")
	}
	status, err := s.statusRepo.GetByName(&constance.Status_Approve)
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
	status, err := s.statusRepo.GetByName(&constance.Status_Review)
	if err != nil {
		return err
	}
	order.Status = *status
	return s.orderRepo.Update(id, order)
}

func (s *OrderService) InsertCustomOrder(order *model.OrderInputRequest) error {
	status, err := s.statusRepo.GetByName(&constance.Status_Pending)
	if err != nil {
		return err
	}
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var orderEntity Entities.Order
	orderEntity.UserID = uuid.MustParse(order.UserID)
	orderEntity.Status = *status
	orderEntity.Items = order.Items
	err = s.orderRepo.Insert(&orderEntity)
	if err != nil {
		tx.Rollback()
		return err
	}
	orderType, err := s.orderTypeRepo.GetByID(&order.OrderTypeID)
	if err != nil {
		tx.Rollback()
		return err
	}
	if orderType == nil {
		tx.Rollback()
		return errors.New("order type not found")
	}
	if orderType.Name == constance.Types_Vow {
		parsedDate, err := time.Parse("2006-01-02", order.Deadline)
		if err != nil {
			return err
		}
		vowRecord := Entities.VowRecord{
			Vow:        order.Vow,
			Deadline:   time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.UTC),
			Note:       order.Note,
			UserID:     uuid.MustParse(order.UserID),
			ServiceID:  uuid.MustParse(order.ServiceID),
			VowOrderID: orderEntity.ID,
		}
		err = s.vowRecordRepo.Insert(&vowRecord)
		if err != nil {
			tx.Rollback()
			return errors.New("failed to insert vow record")
		}
	} else if orderType.Name == constance.Types_Fulfill {
		vowRecord := Entities.VowRecord{
			FulfilledOrderID: orderEntity.ID,
		}
		err = s.vowRecordRepo.Update(&order.VowRecordID, &vowRecord)
		if err != nil {
			tx.Rollback()
			return errors.New("failed to update vow record")
		}
	}
	return nil
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
	order.Price = data.Price
	return s.orderRepo.Update(&data.OrderID, order)
}
