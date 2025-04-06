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
	NotificationUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/notification"
	orderTypeUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/order_type"
	packageUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/package"
	recommendationUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/recommendation"
	serviceUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/service"
	statusUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/status"
	vowRecordUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/vow_record"

	"gorm.io/gorm"
)

type OrderUsecase interface {
	Insert(order *model.OrderInputRequest) (*Entities.Order, error)
	GetAll(config *model.Pagination) ([]*Entities.Order, *model.Pagination, error)
	GetByID(id *string) (*Entities.Order, error)
	GetByStatus(statusID *uuid.UUID, config *model.Pagination) ([]*Entities.Order, *model.Pagination, error)
	Update(id *string, order *Entities.Order) error
	Delete(id *string) error
	Hook(id *string) error
	CancelOrder(id *string, cancelReason *string) error
	SubmitOrder(order *model.SubmitOrderRequest) error
	CompleteOrder(id *string) error
	InsertCustomOrder(order *model.OrderInputRequest) (*Entities.Order, error)
	AcceptOrder(data *model.ConfirmOrderRequest) error
	ApproveOrder(id *string) (*Entities.Order, error)
	GetByUserID(userID *string, config *model.Pagination) ([]*Entities.Order, *model.Pagination, error)
	GetByUserIDAndStatusID(userID *string, statusID *uuid.UUID, config *model.Pagination) ([]*Entities.Order, *model.Pagination, error)
}

type OrderService struct {
	db               *gorm.DB
	orderRepo        OrderRepository
	serviceRepo      serviceUsecase.ServiceRepository
	statusRepo       statusUsecase.StatusRepository
	packageRepo      packageUsecase.PackageRepository
	vowRecordRepo    vowRecordUsecase.VowRecordRepository
	orderTypeRepo    orderTypeUsecase.OrderTypeRepository
	notificationRepo NotificationUsecase.NotificationRepository
	recommendationRepo recommendationUsecase.RecommendationRepository 
}

func NewOrderService(orderRepo OrderRepository, serviceRepo serviceUsecase.ServiceRepository, statusRepo statusUsecase.StatusRepository, db *gorm.DB, packageRepo packageUsecase.PackageRepository, vowRecordRepo vowRecordUsecase.VowRecordRepository, orderTypeRepo orderTypeUsecase.OrderTypeRepository, notificationRepo NotificationUsecase.NotificationRepository, recommendationRepo recommendationUsecase.RecommendationRepository,) OrderUsecase {
	return &OrderService{
		db:               db,
		orderRepo:        orderRepo,
		serviceRepo:      serviceRepo,
		statusRepo:       statusRepo,
		packageRepo:      packageRepo,
		vowRecordRepo:    vowRecordRepo,
		orderTypeRepo:    orderTypeRepo,
		notificationRepo: notificationRepo,
        recommendationRepo: recommendationRepo,
	}
}

func (s *OrderService) Insert(order *model.OrderInputRequest) (*Entities.Order, error) {
    status, err := s.statusRepo.GetByName(&constance.Status_Unpaid)
    if err != nil {
        return nil,err
    }
    packages, err := s.packageRepo.GetByID(&order.PackageID)
    if err != nil {
        return nil,err
    }
    if packages == nil {
        return nil,errors.New("package not found")
    }
	// log.Println("Order Service: ", order)
    // client, err := omise.NewClient(config.OmisePublicKey, config.OmiseSecretKey)
    // if err != nil {
    //     return nil,err
    // }
    // source := &omise.Source{}
    tx := s.db.Begin()
    // defer func() {
    //     if r := recover(); r != nil {
    //         tx.Rollback()
    //     }
    // }()
    // err = client.Do(source, &operations.CreateSource{
    //     Amount:   int64(order.Price * 100),
    //     Currency: "thb",
    //     Type:     "promptpay",
    // })
    // if err != nil {
    //     return nil,err
    // }
    // charge := &omise.Charge{}
    // err = client.Do(charge, &operations.CreateCharge{
    //     Amount:   source.Amount,
    //     Currency: source.Currency,
    //     Source:   source.ID,
    // })
    // if err != nil {
    //     tx.Rollback()
    //     return nil,err
    // }
    var transaction Entities.Transaction
    transaction.Price = order.Price
    // transaction.ChargeID = charge.ID
    // transaction.Charge = *charge
	
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
        return nil,err
    }
    orderType, err := s.orderTypeRepo.GetByID(&order.OrderTypeID)
    if err != nil {
        return nil,err
    }
    if orderType.Name == constance.Types_Vow {
        parsedDate, err := time.Parse("2006-01-02", order.Deadline)
        if err != nil {
            return nil,err
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
            return nil,errors.New("failed to insert vow record")
        }
		
        var orderCount int64
        if err := s.db.Model(&Entities.Order{}).
            Where("user_id = ?", orderEntity.UserID).
            Count(&orderCount).Error; err != nil {
            tx.Rollback()
            return nil,err
        }

		if orderCount > 2 {
			var previousOrder Entities.Order
			if err := s.db.
				Where("user_id = ?", orderEntity.UserID).
				Order("created_at desc").
				Offset(1).
				First(&previousOrder).Error; err == nil {
		
				recEntity := &Entities.Recommendation{
					Current_service_id: previousOrder.ServiceID,
					Next_service_id:    orderEntity.ServiceID,
					Total:              0,
				}
				if err := s.recommendationRepo.Insert(recEntity); err != nil {
					tx.Rollback()
					return nil,err
				}
			}
		}
    } else if orderType.Name == constance.Types_Fulfill {
        vowRecord := Entities.VowRecord{
            FulfilledOrderID: orderEntity.ID,
        }
        err = s.vowRecordRepo.Update(&order.VowRecordID, &vowRecord)
        if err != nil {
            return nil, errors.New("failed to update vow record")
        }
    }
    tx.Commit()
	return &orderEntity, nil
}

func (s *OrderService) Hook(ChargeID *string) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err := s.orderRepo.GetAndUpdateByChargeID(*ChargeID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
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
		tx.Rollback()
		return errors.New("order is already cancelled")
	}
	order.Status = *status
	order.CancellationReason = *cancelReason
	err = s.orderRepo.Update(id, order)
	if err != nil {
		tx.Rollback()
		return err
	}
	s.notificationRepo.Insert(&Entities.Notification{
		Header: "Order Cancelled",
		Body:   "Your order has been cancelled.",
		UserID: order.UserID,
	})
	tx.Commit()
	return nil
}

func (s *OrderService) ApproveOrder(id *string) (*Entities.Order, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if order.Status.Name != constance.Status_Confirm {
		return nil, errors.New("order is not confirm")
	}
	status, err := s.statusRepo.GetByName(&constance.Status_Unpaid)
	if err != nil {
		return nil, err
	}
	order.Status = *status
	client, err := omise.NewClient(config.OmisePublicKey, config.OmiseSecretKey)
	if err != nil {
		return nil, err
	}
	source := &omise.Source{}
	err = client.Do(source, &operations.CreateSource{

		Amount:   int64(order.Price * 100),
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
	transaction.Price = order.Price
	transaction.ChargeID = charge.ID
	transaction.Charge = *charge
	order.Transaction = transaction
	order.Status = *status
	err = s.orderRepo.Update(id, order)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	s.notificationRepo.Insert(&Entities.Notification{
		Header: "Order Approved",
		Body:   "Your order has been approved.",
		UserID: order.UserID,
	})
	tx.Commit()
	return order, nil
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
	err = s.orderRepo.Update(&order.OrderID, orderEntity)
	if err != nil {
		tx.Rollback()
		return err
	}
	s.notificationRepo.Insert(&Entities.Notification{
		Header: "Order Submitted",
		Body:   "Your order has been submitted.",
		UserID: orderEntity.UserID,
	})
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
		tx.Rollback()
		return err
	}
	status, err := s.statusRepo.GetByName(&constance.Status_Review)
	if err != nil {
		tx.Rollback()
		return err
	}
	order.Status = *status
	err = s.orderRepo.Update(id, order)
	if err != nil {
		tx.Rollback()
		return err
	}
	s.notificationRepo.Insert(&Entities.Notification{
		Header: "Order Completed",
		Body:   "Your order has been completed.",
		UserID: order.UserID,
	})
	tx.Commit()
	return nil
}

func (s *OrderService) InsertCustomOrder(order *model.OrderInputRequest) (*Entities.Order, error) {
	status, err := s.statusRepo.GetByName(&constance.Status_Pending)
	if err != nil {
		return nil, err
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
	orderEntity.ServiceID = uuid.MustParse(order.ServiceID)
	err = s.orderRepo.Insert(&orderEntity)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	orderType, err := s.orderTypeRepo.GetByID(&order.OrderTypeID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if orderType == nil {
		tx.Rollback()
		return nil, errors.New("order type not found")
	}
	if orderType.Name == constance.Types_Vow {
		parsedDate, err := time.Parse("2006-01-02", order.Deadline)
		if err != nil {
			return nil, err
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
			return nil, errors.New("failed to insert vow record")
		}
	} else if orderType.Name == constance.Types_Fulfill {
		vowRecord := Entities.VowRecord{
			FulfilledOrderID: orderEntity.ID,
		}
		err = s.vowRecordRepo.Update(&order.VowRecordID, &vowRecord)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("failed to update vow record")
		}
	}
	return &orderEntity, nil
}

func (s *OrderService) AcceptOrder(data *model.ConfirmOrderRequest) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()

		}
	}()
	order, err := s.orderRepo.GetByID(&data.OrderID)
	if err != nil {
		tx.Rollback()
		return err
	}
	if order.Status.Name != constance.Status_Pending {
		tx.Rollback()
		return errors.New("order is not pending")
	}
	status, err := s.statusRepo.GetByName(&constance.Status_Confirm)
	if err != nil {
		tx.Rollback()
		return err
	}
	order.Status = *status
	order.Price = data.Price
	err = s.orderRepo.Update(&data.OrderID, order)
	if err != nil {
		tx.Rollback()
		return err
	}
	s.notificationRepo.Insert(&Entities.Notification{
		Header: "Order Accepted",
		Body:   "Your order has been accepted.",
		UserID: order.UserID,
	})

	tx.Commit()
	return nil
}

func (s *OrderService) GetByUserID(userID *string, config *model.Pagination) ([]*Entities.Order, *model.Pagination, error) {
	if config.PageSize <= 0 {
		config.PageSize = 10
	}
	if config.CurrentPage <= 0 {
		config.CurrentPage = 1
	}
	orders, totalRecords, err := s.orderRepo.GetByUserID(userID, config)
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

func (s *OrderService) GetByUserIDAndStatusID(userID *string, statusID *uuid.UUID, config *model.Pagination) ([]*Entities.Order, *model.Pagination, error) {
	if config.PageSize <= 0 {
		config.PageSize = 10
	}
	if config.CurrentPage <= 0 {
		config.CurrentPage = 1
	}
	orders, totalRecords, err := s.orderRepo.GetByUserIDAndStatusID(userID, statusID, config)
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
