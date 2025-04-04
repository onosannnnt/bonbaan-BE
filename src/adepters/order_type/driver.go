package orderTypeAdapter

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	orderTypeUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/order_type"
	"gorm.io/gorm"
)

type OrderTypeDriver struct {
	db *gorm.DB
}

func NewOrderTypeDriver(db *gorm.DB) orderTypeUsecase.OrderTypeRepository {
	return &OrderTypeDriver{
		db: db,
	}
}

func (d *OrderTypeDriver) Insert(orderType *Entities.OrderType) error {
	if err := d.db.Create(orderType).Error; err != nil {
		return err
	}
	return nil
}

func (d *OrderTypeDriver) GetAll() (*[]Entities.OrderType, error) {
	// Use a slice instead of a pointer to a slice for proper GORM handling.
	var orderTypes []Entities.OrderType
	if err := d.db.Find(&orderTypes).Error; err != nil {
		return nil, err
	}
	return &orderTypes, nil
}

func (d *OrderTypeDriver) GetByID(id *string) (*Entities.OrderType, error) {
	var orderType Entities.OrderType
	if err := d.db.First(&orderType, id).Error; err != nil {
		return nil, err
	}
	return &orderType, nil
}

func (d *OrderTypeDriver) Update(id *string, orderType *Entities.OrderType) error {
	if err := d.db.Model(&Entities.OrderType{}).Where("id = ?", id).Updates(orderType).Error; err != nil {
		return err
	}
	return nil
}

func (d *OrderTypeDriver) Delete(id *string) error {
	if err := d.db.Where("id = ?", id).Delete(&Entities.OrderType{}).Error; err != nil {
		return err
	}
	return nil
}
