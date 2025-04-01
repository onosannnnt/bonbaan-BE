package orderAdepter

import (
	"github.com/google/uuid"
	"github.com/onosannnnt/bonbaan-BE/src/constance"
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	"github.com/onosannnnt/bonbaan-BE/src/model"
	orderUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/order"
	statusUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/status"
	"gorm.io/gorm"
)

type OrderDriver struct {
	db         *gorm.DB
	statusRepo statusUsecase.StatusUsecase
}

func NewOrderDriver(db *gorm.DB, statusRepo statusUsecase.StatusUsecase) orderUsecase.OrderRepository {
	return &OrderDriver{
		db:         db,
		statusRepo: statusRepo,
	}
}

func (d *OrderDriver) GetDefaultStatus() (*Entities.Status, error) {
	var selectStatus Entities.Status
	if err := d.db.Where("name = ?", "unpaid").First(&selectStatus).Error; err != nil {
		return nil, err
	}
	return &selectStatus, nil
}

func (d *OrderDriver) Insert(order *Entities.Order) error {
	if err := d.db.Create(order).Error; err != nil {
		return err
	}
	return nil
}

func (d *OrderDriver) GetAll(config *model.Pagination) ([]*Entities.Order, int64, error) {
	var selectOrder []*Entities.Order
	var totalRecords int64

	if err := d.db.Model(&Entities.Order{}).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	if err := d.db.Preload("Package").Preload("Package.OrderType").Preload("Status").Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Omit("password")
	}).Order("created_at desc").
		Limit(config.PageSize).Offset((config.CurrentPage - 1) * config.PageSize).
		Find(&selectOrder).Error; err != nil {
		return nil, 0, err
	}
	return selectOrder, totalRecords, nil
}

func (d *OrderDriver) GetByID(id *string) (*Entities.Order, error) {
	var selectOrder Entities.Order
	if err := d.db.Preload("Status").Preload("Transaction").Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Omit("password")

	}).Preload("Package").Preload("Package.OrderType").Where("id = ?", id).First(&selectOrder).Error; err != nil {
		return nil, err
	}
	return &selectOrder, nil
}

func (d *OrderDriver) Update(id *string, order *Entities.Order) error {
	if err := d.db.Model(order).Updates(order).Where("id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (d *OrderDriver) Delete(id *string) error {
	if err := d.db.Where("id = ?", *id).Delete(&Entities.Order{}).Error; err != nil {
		return err
	}
	return nil
}

func (d *OrderDriver) GetAndUpdateByChargeID(chargeID string) error {
	var selectOrder Entities.Order
	if err := d.db.Preload("Package").Joins("JOIN transactions ON transactions.id = orders.transaction_id").
		Where("transactions.charge_id = ?", chargeID).
		First(&selectOrder).Error; err != nil {
		return err
	}
	processingOrder, err := d.statusRepo.GetStatusByName(&constance.Status_Processing)
	if err != nil {
		return err
	}
	selectOrder.StatusID = processingOrder.ID
	if err := d.db.Model(&selectOrder).Update("status_id", processingOrder.ID).Error; err != nil {
		return err
	}

	return nil
}

func (d *OrderDriver) GetByStatusID(status *uuid.UUID, config *model.Pagination) ([]*Entities.Order, int64, error) {
	var selectOrder []*Entities.Order
	var totalRecords int64

	if err := d.db.Model(&Entities.Order{}).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}
	if err := d.db.Preload("Status").Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Omit("password")

	}).Preload("Package").Preload("Package.OrderType").Joins("JOIN statuses ON statuses.id = orders.status_id").
		Where("statuses.id = ?", &status).
		Limit(config.PageSize).Offset((config.CurrentPage - 1) * config.PageSize).Find(&selectOrder).Error; err != nil {
		return nil, 0, err
	}
	return selectOrder, totalRecords, nil
}

func (d *OrderDriver) GetByUserID(userID *string, config *model.Pagination) ([]*Entities.Order, int64, error) {
	var selectOrder []*Entities.Order
	var totalRecords int64

	if err := d.db.Model(&Entities.Order{}).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}
	if err := d.db.Preload("Status").Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Omit("password")

	}).Preload("Package").Preload("Package.OrderType").Where("user_id = ?", userID).
		Limit(config.PageSize).Offset((config.CurrentPage - 1) * config.PageSize).Find(&selectOrder).Error; err != nil {
		return nil, 0, err
	}
	return selectOrder, totalRecords, nil
}

func (d *OrderDriver) GetByUserIDAndStatusID(userID *string, statusID *uuid.UUID, config *model.Pagination) ([]*Entities.Order, int64, error) {
	var selectOrder []*Entities.Order
	var totalRecords int64

	if err := d.db.Model(&Entities.Order{}).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}
	if err := d.db.Preload("Status").Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Omit("password")

	}).Preload("Package").Preload("Package.OrderType").Joins("JOIN statuses ON statuses.id = orders.status_id").
		Where("user_id = ? AND statuses.id = ?", userID, statusID).
		Limit(config.PageSize).Offset((config.CurrentPage - 1) * config.PageSize).Find(&selectOrder).Error; err != nil {
		return nil, 0, err
	}
	return selectOrder, totalRecords, nil
}
